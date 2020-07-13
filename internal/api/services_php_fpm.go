package api

import (
	"encoding/json"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
)

func (s *server) handlePhpFpmService() http.HandlerFunc {
	type request struct {
		Version string `json:"version"`
		Action  string `json:"action"`
	}

	rules := govalidator.MapData{
		"version": []string{"required", "in:7.4,7.3,7.2"},
		"action":  []string{"required", "in:stop,start,restart"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		resp := Response{}

		s.validateRequest(w, r, rules, req)

		out, err := s.service.Run("service", []string{"php" + req.Version + "-fpm", req.Action})
		if err != nil {
			s.logger.Println(err, ":", string(out), "\n request: ", req)

			jsonError(w, err, http.StatusBadRequest)

			return
		}

		resp.Success = true

		switch req.Action {
		case "stop":
			resp.Output = "stopped php" + req.Version + "-fpm"
		case "start":
			resp.Output = "started php" + req.Version + "-fpm"
		default:
			resp.Output = "restarted php" + req.Version + "-fpm"
		}

		js, err := json.Marshal(resp)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(js); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) validateRequest(w http.ResponseWriter, r *http.Request, rules govalidator.MapData, data interface{}) {
	errors := govalidator.New(govalidator.Options{
		Request:         r,
		Rules:           rules,
		Data:            &data,
		RequiredDefault: true,
	}).ValidateJSON()

	if len(errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := map[string]interface{}{"errors": errors}
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}
