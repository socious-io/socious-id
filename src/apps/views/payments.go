package views

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"github.com/stripe/stripe-go/v81"
)

func paymentsGroup(router *gin.Engine) {
	g := router.Group("payments")

	g.GET("/fiat/cards", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		stripeCustomerID := identity.Meta["stripe_customer_id"]
		if stripeCustomerID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID has not been set on this identity"})
			return
		}

		fiatService := config.Config.Payment.Fiats[0]
		paymentMethods, err := fiatService.FetchCards(stripeCustomerID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"cards": paymentMethods,
		})
	})

	g.POST("/fiat/cards", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)

		stripeCustomerID, email := identity.Meta["stripe_customer_id"], identity.Meta["email"].(string)
		fiatService := config.Config.Payment.Fiats[0]

		form := new(AddCardForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var (
			customer       *stripe.Customer
			paymentMethod  *stripe.PaymentMethod
			identityEntity interface{}
			err            error
		)
		if form.Token == nil && stripeCustomerID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment source card could not be found"})
			return
		} else if stripeCustomerID == nil {
			customer, err = fiatService.AddCustomer(email)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			switch identity.Type {
			case models.IdentityTypeUsers:
				identityEntity, err = models.GetUser(identity.ID)
			case models.IdentityTypeOrganizations:
				identityEntity, err = models.GetOrganization(identity.ID)
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			switch v := identityEntity.(type) {
			case *models.User:
				v.StripeCustomerID = &customer.ID
				err = v.UpdateProfile(ctx)
			case *models.Organization:
				v.StripeCustomerID = &customer.ID
				err = v.Update(ctx)
			}

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			stripeCustomerID = customer.ID
		}

		paymentMethod, err = fiatService.AttachPaymentMethod(stripeCustomerID.(string), *form.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, paymentMethod)
	})

	g.DELETE("/fiat/cards/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := c.Param("id")

		stripeCustomerID := identity.Meta["stripe_customer_id"]
		if stripeCustomerID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID has not been set on this identity"})
			return
		}

		fiatService := config.Config.Payment.Fiats[0]
		err := fiatService.DeleteCard(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.GET("/fiat/payout", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		oc, err := models.GetOauthConnectByIdentityId(identity.ID, models.OauthConnectedProvidersStripeJp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fiatService := config.Config.Payment.Fiats[0]
		acc, err := fiatService.FetchAccount(oc.MatrixUniqueID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, json.RawMessage(acc.LastResponse.RawJSON))
	})

	g.GET("/fiat/payout/connect", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		userRedirectUrl := c.Query("redirect_url")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		fiatService := config.Config.Payment.Fiats[0]
		account, err := fiatService.CreateAccount("JP")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		redirectURL := fmt.Sprintf("%s?stripe_account=%s", fiatService.Callback, account.ID)

		link, err := fiatService.CreateAccountLink(account, redirectURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		oc, err := models.GetOauthConnectByIdentityId(identity.ID, models.OauthConnectedProvidersStripeJp)
		if err != nil && oc == nil {
			oc = &models.OauthConnect{
				IdentityID:     identity.ID,
				MatrixUniqueID: account.ID,
				AccessToken:    "",
				Status:         models.UserStatusTypeInactive,
				Provider:       models.OauthConnectedProvidersStripeJp, //WARNING: Hardcoded
				RedirectURL:    &userRedirectUrl,
				IsConfirmed:    false,
			}

			if err := oc.Upsert(ctx); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"url": link,
		})
	})

	g.GET("/fiat/payout/callback/stripe", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		stripeAccount := c.Query("stripe_account")

		oc, err := models.GetOauthConnectByMUI(stripeAccount, models.OauthConnectedProvidersStripeJp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if oc.IsConfirmed {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?status=success", *oc.RedirectURL))
			return
		}

		if stripeAccount == "" {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?status=failed&error=%s",
				*oc.RedirectURL,
				"stripe account query is empty",
			))
			return
		}

		fiatService := config.Config.Payment.Fiats[0]

		acc, _ := fiatService.FetchAccount(stripeAccount)
		accountJson := types.JSONText(acc.LastResponse.RawJSON)

		//Updating the OauthConnect
		oc.Meta = &accountJson
		oc.Status = models.UserStatusTypeInactive
		oc.IsConfirmed = true
		if acc.PayoutsEnabled {
			oc.Status = models.UserStatusTypeActive
		}

		err = oc.Upsert(ctx)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?status=failed&error=%s",
				*oc.RedirectURL,
				err.Error(),
			))
		}

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?status=success", *oc.RedirectURL))
	})

	g.GET("/crypto/wallets", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		wallets, err := models.GetWallets(identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, wallets)
	})

	g.POST("/crypto/wallets", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(AddWalletForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		wallet := new(models.Wallet)
		utils.Copy(form, wallet)
		wallet.IdentityID = identity.ID

		if err := wallet.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, wallet)
	})
}
