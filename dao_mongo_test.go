package authtokensvc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/joelywz/authtokensvc"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDao(t *testing.T) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017"))

	if err != nil {
		t.Fatalf("failed initialize client, %v", err)
	}

	err = client.Connect(context.Background())

	if err != nil {
		t.Fatalf("failed connect to database, %v", err)
	}

	db := client.Database("test_auth_service")

	db.Drop(context.Background())

	collection := db.Collection("auth_tokens")

	dao := authtokensvc.NewMongoDao(collection)

	newToken := authtokensvc.Token{
		RefreshID:        gonanoid.Must(32),
		AccessID:         gonanoid.Must(32),
		RefreshExpiresAt: time.Now(),
		AccessExpiresAt:  time.Now(),
		UserID:           "2",
	}

	t.Run("Save", func(t *testing.T) {
		if err := dao.SaveToken(&newToken); err != nil {
			t.Fatalf("failed to save token, %v", err)
		}
	})

	t.Run("GetAccessToken", func(t *testing.T) {
		retrivedToken, err := dao.GetAccessToken(newToken.AccessID)

		if err != nil {
			t.Fatalf("failed to get access token, %v", err)
		}

		if err := compareTokens(retrivedToken, &newToken); err != nil {
			t.Fatalf("expected token data to be identical, %v", err)
		}

	})

	t.Run("GetRefreshToken", func(t *testing.T) {
		retrivedToken, err := dao.GetRefreshToken(newToken.RefreshID)

		if err != nil {
			t.Fatalf("failed to get access token, %v", err)
		}

		if err := compareTokens(retrivedToken, &newToken); err != nil {
			t.Fatalf("expected token data to be identical, %v", err)
		}
	})

	t.Run("DeleteToken1", func(t *testing.T) {
		err := dao.Delete(newToken.AccessID)

		if err != nil {
			t.Fatalf("failed to delete token, %v", err)
		}

		token, _ := dao.GetAccessToken(newToken.AccessID)

		if token != nil {
			t.Fatalf("expected access token to be deleted, but token is still located in database")
		}
	})

	t.Run("DeleteToken2", func(t *testing.T) {
		err := dao.Delete(newToken.RefreshID)

		if err != nil {
			t.Fatalf("failed to delete token, %v", err)
		}

		token, _ := dao.GetAccessToken(newToken.RefreshID)

		if token != nil {
			t.Fatalf("expected access token to be deleted, but token is still exist in database")
		}

		if err := dao.SaveToken(&newToken); err != nil {
			t.Fatalf("failed to save token, %v", err)
		}

	})

	t.Run("DeleteAll", func(t *testing.T) {
		newToken2 := authtokensvc.Token{
			RefreshID:        gonanoid.Must(32),
			AccessID:         gonanoid.Must(32),
			RefreshExpiresAt: time.Now(),
			AccessExpiresAt:  time.Now(),
			UserID:           "2",
		}

		newToken3 := authtokensvc.Token{
			RefreshID:        gonanoid.Must(32),
			AccessID:         gonanoid.Must(32),
			RefreshExpiresAt: time.Now(),
			AccessExpiresAt:  time.Now(),
			UserID:           "1",
		}

		dao.SaveToken(&newToken2)
		dao.SaveToken(&newToken3)

		err := dao.DeleteAll("2")

		if err != nil {
			t.Fatalf("failed to delete all tokens, %v", err)
		}

		token, _ := dao.GetAccessToken(newToken3.AccessID)

		if token == nil {
			t.Fatalf("expoected tokens of the user ID to be deleted, but token of other user ID was deleted")
		}

		token, _ = dao.GetAccessToken(newToken2.AccessID)

		if token != nil {
			t.Fatalf("expoected token of the user ID to be deleted, but still exist in database")
		}

		token, _ = dao.GetAccessToken(newToken.AccessID)

		if token != nil {
			t.Fatalf("expoected token of the user ID to be deleted, but still exist in database")
		}

	})

}

func compareTokens(a *authtokensvc.Token, b *authtokensvc.Token) error {
	if a.AccessID != b.AccessID {
		return errors.New("access token id does not match")
	}

	if a.RefreshID != b.RefreshID {
		return errors.New("refresh token id does not match")
	}

	if !a.RefreshExpiresAt.Truncate(time.Millisecond).Equal(b.RefreshExpiresAt.Truncate(time.Millisecond)) {
		return errors.New("refresh token expiry does not match")
	}

	if !a.AccessExpiresAt.Truncate(time.Millisecond).Equal(b.AccessExpiresAt.Truncate(time.Millisecond)) {
		return errors.New("access token expiry does not match")
	}

	return nil
}
