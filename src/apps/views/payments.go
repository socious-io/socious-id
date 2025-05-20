package views

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/card"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentmethod"
	"github.com/stripe/stripe-go/v82/paymentsource"
)

func createCard(token string, email string) (*stripe.Customer, *stripe.PaymentSource, error) {
	paymentMethod, err := paymentmethod.New(&stripe.PaymentMethodParams{
		Type: stripe.String(stripe.PaymentMethodTypeCard),
		Card: &stripe.PaymentMethodCardParams{
			Token: stripe.String(token),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	customer, err := customer.New(&stripe.CustomerParams{
		Email:         stripe.String(email),
		PaymentMethod: stripe.String(paymentMethod.ID),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethod.ID),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	paymentSource, err := paymentsource.New(&stripe.PaymentSourceParams{
		Source: &stripe.PaymentSourceSourceParams{
			Token: stripe.String(token),
		},
		Customer: stripe.String(customer.ID),
	})
	if err != nil {
		return nil, nil, err
	}

	return customer, paymentSource, nil
}

func deleteCard(customerID string, cardID string) (*stripe.Card, error) {
	result, err := card.Del(cardID, &stripe.CardParams{Customer: stripe.String(customerID)})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func paymentsGroup(router *gin.Engine) {
	g := router.Group("payments")

	g.GET("cards", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		paginate := c.MustGet("paginate").(database.Paginate)
		identity := c.MustGet("identity").(*models.Identity)

		page, limit := c.MustGet("page"), c.MustGet("limit")

		cards, total, err := models.GetCards(identity.ID, paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":    page,
			"limit":   limit,
			"results": cards,
			"total":   total,
		})
	})

	g.GET("cards/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := uuid.MustParse(c.Param("id"))

		card, err := models.GetCard(id, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusFound, card)
	})

	g.POST("/cards", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(AddCardForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		scustomer, scard, err := createCard(form.Token, identity.MetaMap["email"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to create card: %w", err)})
			return
		}

		card := &models.Card{
			IdentityID: identity.ID,
			Customer:   scustomer.ID,
			Card:       scard.ID,
		}

		if err := card.CreateCard(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, card)
	})

	g.PUT("/cards/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)
		id := c.MustGet("id").(uuid.UUID)

		form := new(AddCardForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		card, err := models.GetCard(id, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("couldn't find the card to update: %w", err)})
			return
		}

		scustomer, scard, err := createCard(form.Token, identity.MetaMap["email"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to create card: %w", err)})
			return
		}
		card.Customer = scustomer.ID
		card.Card = scard.ID

		if err := card.UpdateCard(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, card)
	})

	g.DELETE("/cards/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)
		id := c.MustGet("id").(uuid.UUID)

		card, err := models.GetCard(id, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("couldn't find the card to delete: %w", err)})
			return
		}

		_, err = deleteCard(card.Customer, card.Card)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("couldn't delete card: %w", err)})
			return
		}

		if err := card.DeleteCard(ctx); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "success"})
	})

	g.GET("/wallets", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		wallets, err := models.GetWallets(identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusFound, gin.H{"wallets": wallets})
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

		c.JSON(http.StatusBadRequest, gin.H{"wallet": wallet})
	})
}
