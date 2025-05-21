package views

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/config"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
)

func paymentsGroup(router *gin.Engine) {
	g := router.Group("payments")

	g.GET("cards", auth.LoginRequired(), func(c *gin.Context) {
		// paginate := c.MustGet("paginate").(database.Paginate)
		identity := c.MustGet("identity").(*models.Identity)

		// page, limit := c.MustGet("page"), c.MustGet("limit")
		fmt.Println("identity", identity)
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

	g.POST("/cards", auth.LoginRequired(), func(c *gin.Context) {
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
		}

		paymentMethod, err = fiatService.AttachPaymentMethod(customer.ID, *form.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"card":     paymentMethod,
			"customer": customer,
			"identity": identity,
		})
	})

	g.DELETE("/cards/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := c.MustGet("id").(string)

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

	g.GET("/wallets", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		wallets, err := models.GetWallets(identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, wallets)
	})

	g.POST("/wallets", auth.LoginRequired(), func(c *gin.Context) {
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
