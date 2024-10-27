package auth

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	telegramloginwidget "github.com/LipsarHQ/go-telegram-login-widget"
	"github.com/RGITHackathonFall2024/auth/internal/consts"
	"github.com/RGITHackathonFall2024/auth/internal/logic/user"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) (string, error) {
	jwtSecret := []byte(os.Getenv(consts.EnvJwtSecret))

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"sub": userID,
	}

	return token.SignedString(jwtSecret)
}

func VerifyHash(log *slog.Logger, tgAuthData *telegramloginwidget.AuthorizationData) error {
	log = log.WithGroup("verify-hash")

	if time.Unix(tgAuthData.AuthDate, 0).Add(time.Hour * 24).Before(time.Now()) {
		return &ErrStaleAuthData{}
	}

	botToken := os.Getenv(consts.EnvTgToken)

	log.Debug("Bot token", slog.String("bot_token", botToken))

	/*
		secret := sha256.Sum256([]byte(botToken))
		log.Debug("Secret", slog.String("secret", string(secret[:])))

		dataCheckString := strings.Join(
			slices.Sorted(utils.Map(
				slices.Values([][]string{
					{"id", strconv.FormatInt(id, 10)},
					{"first_name", firstName},
					{"last_name", lastName},
					{"username", username},
					{"photo_url", photoUrl},
					{"auth_date", strconv.FormatUint(authDate, 10)},
				}),
				func(field []string) string { return strings.Join(field, "=") },
			)),
			"\n",
		)
		log.Debug("Data check string", slog.String("data_check_string", dataCheckString))

		gotHash := hmac.New(sha256.New, secret[:])
		if _, err := gotHash.Write([]byte(dataCheckString)); err != nil {
			log.Error("Error hashing data check string", slog.String("err", err.Error()))
			return err
		}

		gotHashHex := make([]byte, hex.EncodedLen(gotHash.Size()))
		if n := hex.Encode(gotHashHex, gotHash.Sum(nil)); n != len(gotHashHex) {
			log.Error("Error encoding hash", slog.Int("n", n), slog.Int("expected", len(gotHashHex)))
			return &ErrEncodingHash{}
		}

		log.Debug("Got hash", slog.String("got_hash", string(gotHashHex)))

		if !hmac.Equal([]byte(hash), gotHashHex) {
			return &ErrInvalidHash{}
		}
	*/

	if err := tgAuthData.Check(botToken); err != nil {
		return &ErrInvalidHash{}
	}

	return nil
}

func GetToken(log *slog.Logger, c *fiber.Ctx) (string, error) {
	log = log.WithGroup("get-token")

	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		log.Error("Missing token", slog.String("auth_header", authHeader))
		return "", &ErrMissingToken{}
	}

	return token, nil
}

func GetUserByToken(log *slog.Logger, s *server.Server, token string) (*user.User, error) {
	var usr *user.User
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Unexpected signing method", slog.String("method", t.Method.Alg()))
			return nil, &ErrInvalidToken{}
		}

		exp, err := t.Claims.GetExpirationTime()
		if err != nil {
			log.Error("Error getting expiration time", slog.String("err", err.Error()))
			return nil, &ErrInvalidToken{}
		}

		if exp.Before(time.Now()) {
			log.Error("Token expired", slog.Time("exp", exp.Time))
			return nil, &ErrInvalidToken{}
		}

		idStr, err := t.Claims.GetSubject()
		if err != nil {
			log.Error("Error getting subject", slog.String("err", err.Error()))
			return nil, &ErrInvalidToken{}
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("Error parsing subject", slog.String("err", err.Error()), slog.String("subject", idStr))
			return nil, &ErrInvalidToken{}
		}

		usr, err = user.ByID(s, id)
		if err != nil {
			log.Error("Error getting user", slog.String("err", err.Error()), slog.Int64("telegram_id", id))
			return nil, &ErrInvalidToken{}
		}

		return []byte(os.Getenv(consts.EnvJwtSecret)), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			log.Error("Signature invalid", slog.String("err", err.Error()))
			return nil, &ErrInvalidToken{}
		}

		return nil, err
	}

	return usr, nil
}

func GetUserFromContext(log *slog.Logger, c *fiber.Ctx) (*user.User, error) {
	log = log.WithGroup("get-user-from-context")

	s := server.FromContext(c)
	if s == nil {
		return nil, &server.ErrNoServerInContext{}
	}

	tokenStr, err := GetToken(log, c)
	if err != nil {
		return nil, err
	}

	return GetUserByToken(log, s, tokenStr)
}
