package internal

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	Servicer
}

func NewHandler(hand Servicer) *Handler {
	if hand == nil {
		return nil
	}
	return &Handler{Servicer: hand}

}

// /validation function
func (h *Handler) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	//validate.RegisterValidation()
	return validate.Struct(h)
}

/*func validatepassword(password string) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	matched, _ := regexp.MatchString(pattern, password)
	return matched
}*/

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string      `json:"username" validate:"required,min=5,max=20"`
		Email     string      `json:"email" validate:"required,email"`
		Password  string      `json:"password" validate:"required"`
		Firstname pgtype.Text `json:"firstname" validate:"required"`
		Lastname  pgtype.Text `json:"lastname" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tok, err := h.Servicer.SignUp(r.Context(), req.Username, req.Email, req.Password, req.Firstname, req.Lastname)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		//log.Fatalf("there was error creating user :%v", err)
		return
	}

	resp := map[string]string{tok: "token"}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tok, err := h.Servicer.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//cookies are generated here
	cookie := http.Cookie{
		Name:     "user_session",
		Value:    tok,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	//w.Write([]byte("cookie set!"))
	//fmt.Println(user.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"msg":   "Login successful",
		"email": "req.Email",
	})

	/*resp := map[string]string{tok: "token"}
	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)*/

}
