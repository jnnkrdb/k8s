package healthz

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jnnkrdb/httprdb"
)

var (
	_readyness int = 500
	_liveness  int = 500
)

// complete http.Route object for the api
var (
	LivenessHandler  = httprdb.Route{Request: "GET", SubPath: "/healthz/liveness", Handler: Liveness}
	ReadynessHandler = httprdb.Route{Request: "GET", SubPath: "/healthz/readyness", Handler: Readyness}
)

// receive the current liveness value
func GetState_Liveness() int {
	return _liveness
}

// receive the current readyness value
func GetState_Readyness() int {
	return _readyness
}

// receive the current liveness value
//
// Parameters:
//   - `state` : int > update the value of the liveness state
func SetState_Liveness(state int) {
	_liveness = state
}

// receive the current readyness value
//
// Parameters:
//   - `state` : int > update the value of the readyness state
func SetState_Readyness(state int) {
	_readyness = state
}

// --------------------------------------------------

// receive the handlerfunction for the httprdb.Route
// NOTE: based on gin router
func Liveness(ctx *gin.Context) {
	ctx.IndentedJSON(_liveness, "Status: "+strconv.Itoa(_liveness))
}

// receive the handlerfunction for the httprdb.Route
// NOTE: based on gin router
func Readyness(ctx *gin.Context) {
	ctx.IndentedJSON(_readyness, "Status: "+strconv.Itoa(_readyness))
}
