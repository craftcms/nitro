package api

import (
	"encoding/json"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
)

func (s *server) handleNginxService() http.HandlerFunc {
	type request struct {
		Action string `json:"action"`
	}

	rules := govalidator.MapData{
		"action": []string{"required", "in:stop,start,restart"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		resp := Response{}

		errors := govalidator.New(govalidator.Options{
			Request:         r,
			Rules:           rules,
			Data:            &req,
			RequiredDefault: true,
		}).ValidateJSON()

		// TODO abstract the error check
		if len(errors) > 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := map[string]interface{}{"errors": errors}
			_ = json.NewEncoder(w).Encode(err)
			return
		}

		out, err := s.service.Run("service", []string{"nginx", req.Action})
		if err != nil {
			s.logger.Println(err, ":", string(out), "\n request: ", req)

			jsonError(w, err, http.StatusBadRequest)

			return
		}

		resp.Success = true

		switch req.Action {
		case "stop":
			resp.Output = "stopped nginx"
		case "start":
			resp.Output = "started nginx"
		default:
			resp.Output = "restarted nginx"
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
