package validators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validators "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"users-api/pkg/apperror"
)

type Validator struct {
	validator *validators.Validate
	trans     ut.Translator
}

func NewValidator() (*Validator, error) {
	v := validators.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		validator: v,
		trans:     trans,
	}, nil
}

func (vld *Validator) ValidateStruct(v interface{}) error {
	err := vld.validator.Struct(v)
	if err == nil {
		return nil
	}

	validatorErrs := err.(validators.ValidationErrors)

	var errs = make([]string, len(validatorErrs))
	for i, e := range validatorErrs {
		err = vld.determineError(e)
		_, ok := apperror.IsAppError(err)
		if !ok {
			err = apperror.NewInvalidRequestWithMsg(e.Translate(vld.trans))
		}

		errs[i] = err.Error()
	}

	errMsg := combineErrorMsgs(errs)
	return apperror.NewInvalidRequestWithMsg(errMsg)
}

func (vld *Validator) determineError(e validators.FieldError) error {
	setting := map[string]error{}

	err, ok := setting[fmt.Sprintf("%v.%v", e.Field(), e.Tag())]
	if ok {
		return err
	}

	err, ok = setting[e.Tag()]
	if ok {
		return err
	}

	return e
}

// This function combines error messages following the format from github.com/go-ozzo/ozzo-validation/v4, which is the old validation library we used in trading-orderbook-api.
func combineErrorMsgs(errMsgs []string) string {
	if len(errMsgs) == 0 {
		return ""
	}

	if len(errMsgs) == 1 {
		return errMsgs[0]
	}

	sort.Strings(errMsgs)
	// Remove dot(.) from err messages except the last one.
	for i, msg := range errMsgs[:len(errMsgs)-1] {
		errMsgs[i] = strings.TrimSuffix(msg, ".")
	}

	combinedErrMsg := strings.Join(errMsgs, "; ")
	return combinedErrMsg
}
