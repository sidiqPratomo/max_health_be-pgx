package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/sidiqPratomo/max-health-backend/appvalidator"
	"github.com/sidiqPratomo/max-health-backend/config"
	"github.com/sidiqPratomo/max-health-backend/handler"
	"github.com/sidiqPratomo/max-health-backend/middleware"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type routerOpts struct {
	Ping               *handler.PingHandler
	Authentication     *handler.AuthenticationHandler
	User               *handler.UserHandler
	Doctor             *handler.DoctorHandler
	UserAddress        *handler.UserAddressHandler
	Partner            *handler.PartnerHandler
	Address            *handler.AddressHandler
	Drug               *handler.DrugHandler
	DrugForm           *handler.DrugFormHandler
	DrugClassification *handler.DrugClassificationHandler
	Category           *handler.CategoryHandler
	Cart               *handler.CartHandler
	Telemedicine       *handler.TelemedicineHandler
	Order              *handler.OrderHandler
	Pharmacy           *handler.PharmacyHandler
	OrderPharmacy      *handler.OrderPharmacyHandler
	Report             *handler.ReportHandler
	Stock              *handler.StockHandler
}

type utilOpts struct {
	JwtHelper util.TokenAuthentication
}

func newRouter(h routerOpts, u utilOpts, config *config.Config, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	corsConfig := cors.DefaultConfig()

	router.ContextWithFallback = true

	appvalidator.AppValidator()

	router.Use(
		middleware.Logger(log),
		middleware.RequestIdHandlerMiddleware,
		middleware.ErrorHandlerMiddleware,
		gin.Recovery(),
	)

	authMiddleware := middleware.AuthMiddleware(u.JwtHelper, config)

	userAuthorizationMiddleware := middleware.UserAuthorizationMiddleware
	doctorAuthorizationMiddleware := middleware.DoctorAuthorizationMiddleware
	pharmacyManagerAuthorizationMiddleware := middleware.PharmacyManagerAuthorizationMiddleware
	adminAuthorizationMiddleware := middleware.AdminAuthorizationMiddleware

	corsRouting(router, corsConfig)
	router.NoRoute(handler.NotFoundHandler)
	authenticationRouting(router, h.Authentication)
	addressRouting(router, h.Address, authMiddleware)
	doctorRouting(router, h.Doctor, authMiddleware, doctorAuthorizationMiddleware)
	userRouting(router, h.User, h.Cart, authMiddleware, userAuthorizationMiddleware)
	userAddressRouting(router, h.UserAddress, authMiddleware, userAuthorizationMiddleware)
	partnerRouting(router, h.Partner, authMiddleware, adminAuthorizationMiddleware)
	drugRouting(router, h.Drug, authMiddleware, adminAuthorizationMiddleware, pharmacyManagerAuthorizationMiddleware)
	drugFormRouting(router, h.DrugForm)
	drugClassificationRouting(router, h.DrugClassification)
	pharmacyRouting(router, h.Pharmacy, authMiddleware, pharmacyManagerAuthorizationMiddleware, adminAuthorizationMiddleware)
	categoryRouting(router, h.Category, authMiddleware, adminAuthorizationMiddleware)
	cartRouting(router, h.Cart, authMiddleware, userAuthorizationMiddleware)
	telemedicineRouting(router, h.Telemedicine, authMiddleware, userAuthorizationMiddleware, doctorAuthorizationMiddleware)
	orderRouting(router, h.Order, authMiddleware, userAuthorizationMiddleware, adminAuthorizationMiddleware, pharmacyManagerAuthorizationMiddleware)
	orderPharmacyRouting(router, h.OrderPharmacy, authMiddleware, pharmacyManagerAuthorizationMiddleware, userAuthorizationMiddleware, adminAuthorizationMiddleware)
	reportRouting(router, h.Report, authMiddleware, pharmacyManagerAuthorizationMiddleware, adminAuthorizationMiddleware)
	stockRouting(router, h.Stock, authMiddleware, pharmacyManagerAuthorizationMiddleware)
	pingRouting(router, h.Ping, authMiddleware, userAuthorizationMiddleware, doctorAuthorizationMiddleware, pharmacyManagerAuthorizationMiddleware, adminAuthorizationMiddleware)
	pprofRouting(router)

	return router
}

func userAddressRouting(router *gin.Engine, handler *handler.UserAddressHandler, authMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc) {
	router.POST("/address", authMiddleware, userAuthorizationMiddleware, handler.AddUserAddress)
	router.POST("/address/autofill", authMiddleware, userAuthorizationMiddleware, handler.AddUserAddressAutofill)
	router.PUT("/address/:address_id", authMiddleware, userAuthorizationMiddleware, handler.UpdateUserAddress)
	router.GET("/address", authMiddleware, userAuthorizationMiddleware, handler.GetAllUserAddress)
	router.DELETE("/address/:address_id", authMiddleware, userAuthorizationMiddleware, handler.DeleteUserAddress)
}

func drugRouting(router *gin.Engine, handler *handler.DrugHandler, authMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc, pharmacyManagerAuthorizationMiddleware gin.HandlerFunc) {
	router.GET("/drugs/:drug_id", handler.GetPharmacyDrugByDrugId)
	router.GET("/drugs", handler.GetAllDrugsForListing)

	router.GET("/admin/drugs", authMiddleware, handler.GetAllDrugs)
	router.GET("/admin/drugs/:drug_id", authMiddleware, adminAuthorizationMiddleware, handler.GetDrugByDrugId)
	router.PUT("/admin/drugs/:drug_id", authMiddleware, adminAuthorizationMiddleware, handler.UpdateOneDrug)
	router.POST("/admin/drugs", authMiddleware, adminAuthorizationMiddleware, handler.CreateOneDrug)
	router.DELETE("/admin/drugs/:drug_id", authMiddleware, adminAuthorizationMiddleware, handler.DeleteOneDrug)

	router.GET("/managers/pharmacies/:pharmacy_id/drugs", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetDrugsByPharmacyId)
	router.PATCH("/managers/pharmacies/drugs/:pharmacy_drug_id", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.UpdateDrugsByPharmacyDrugId)
	router.DELETE("/managers/pharmacies/drugs/:pharmacy_drug_id", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.DeleteDrugsByPharmacyDrugId)
	router.POST("/managers/pharmacies/drugs", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.AddDrugsByPharmacyManager)
	router.GET("/managers/pharmacies/drugs/:pharmacy_drug_id/mutation", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetPossibleStockMutation)
	router.POST("/managers/pharmacies/drugs/:pharmacy_drug_id/mutation", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.PostStockMutation)
}

func drugFormRouting(router *gin.Engine, handler *handler.DrugFormHandler) {
	router.GET("/drugs/forms", handler.GetAllDrugForm)
}

func drugClassificationRouting(router *gin.Engine, handler *handler.DrugClassificationHandler) {
	router.GET("/drugs/classifications", handler.GetAllDrugClassification)
}

func userRouting(router *gin.Engine, handler *handler.UserHandler, cartHandler *handler.CartHandler, authMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc) {
	userRouter := router.Group("/users")
	userRouter.PATCH("/profile", authMiddleware, userAuthorizationMiddleware, handler.UpdateData)
	userRouter.GET("/profile", authMiddleware, userAuthorizationMiddleware, handler.GetProfile)
}

func doctorRouting(router *gin.Engine, handler *handler.DoctorHandler, authMiddleware gin.HandlerFunc, doctorAuthorizationMiddleware gin.HandlerFunc) {
	doctorRouter := router.Group("/doctors")
	doctorRouter.PATCH("/profile", authMiddleware, doctorAuthorizationMiddleware, handler.UpdateData)
	doctorRouter.GET("/", handler.GetAllDoctors)
	doctorRouter.GET("/specializations", handler.GetAllDoctorSpecialization)
	doctorRouter.GET("/profile", authMiddleware, doctorAuthorizationMiddleware, handler.GetProfile)
	doctorRouter.GET(":doctor_id", handler.GetProfileForPublic)
	doctorRouter.PATCH("/availability", authMiddleware, doctorAuthorizationMiddleware, handler.UpdateDoctorStatus)
	doctorRouter.GET("/availability", authMiddleware, doctorAuthorizationMiddleware, handler.GetDoctorIsOnline)
}

func partnerRouting(router *gin.Engine, handler *handler.PartnerHandler, authMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc) {
	partnerRouter := router.Group("/partners")

	partnerRouter.POST("", authMiddleware, adminAuthorizationMiddleware, handler.AddPartner)
	partnerRouter.POST("/access-details", authMiddleware, adminAuthorizationMiddleware, handler.SendCredentialsEmail)
	partnerRouter.GET("", authMiddleware, adminAuthorizationMiddleware, handler.GetAllPartners)
	partnerRouter.PATCH("/:pharmacy_manager_id", authMiddleware, adminAuthorizationMiddleware, handler.UpdatePartner)
	partnerRouter.DELETE("/:pharmacy_manager_id", authMiddleware, adminAuthorizationMiddleware, handler.DeletePartner)
}

func pharmacyRouting(router *gin.Engine, handler *handler.PharmacyHandler, authMiddleware gin.HandlerFunc, pharmacyManagerAuthorizationMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc) {
	router.GET("/managers/pharmacies", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetPharmacyByManagerId)
	router.PUT("/pharmacies/:pharmacy_id", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.UpdateOnePharmacy)
	router.DELETE("/pharmacies/:pharmacy_id", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.DeleteOnePharmacy)
	router.POST("/pharmacies", authMiddleware, adminAuthorizationMiddleware, handler.CreateOnePharmacy)
	router.GET("/admin/manager/:pharmacy_manager_id/pharmacies", authMiddleware, adminAuthorizationMiddleware, handler.AdminGetPharmacyByManagerId)
}

func stockRouting(router *gin.Engine, handler *handler.StockHandler, authMiddleware gin.HandlerFunc, pharmacyManagerAuthorizationMiddleware gin.HandlerFunc) {
	router.GET("/managers/stock-change", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetAllStockChanges)
}

func addressRouting(router *gin.Engine, handler *handler.AddressHandler, authMiddleware gin.HandlerFunc) {
	router.GET("/provinces", handler.GetAllProvinces)
	router.GET("/cities", handler.GetAllCitiesByProvinceCode)
	router.GET("/districts", handler.GetAllDistrictsByCityCode)
	router.GET("/subdistricts", handler.GetAllSubdistrictsByDistrictCode)
}

func telemedicineRouting(router *gin.Engine, handler *handler.TelemedicineHandler, authMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc, doctorAuthorizationMiddleware gin.HandlerFunc) {
	router.POST("/chat-rooms", authMiddleware, userAuthorizationMiddleware, handler.UserCreateRoom)
	router.PATCH("/chat-rooms", authMiddleware, doctorAuthorizationMiddleware, handler.DoctorJoinRoom)
	router.POST("/chat-rooms/chats", authMiddleware, handler.PostOneMessage)
	router.GET("/chat-rooms/chats/:room_id", authMiddleware, handler.Listen)
	router.GET("/chat-rooms/:room_id", authMiddleware, handler.GetAllChat)
	router.GET("/chat-rooms", authMiddleware, handler.GetAllChatRoomPreview)
	router.GET("/chat-rooms/requests", authMiddleware, doctorAuthorizationMiddleware, handler.DoctorGetChatRequest)
	router.PATCH("/chat-rooms/:room_id/close-room", authMiddleware, userAuthorizationMiddleware, handler.CloseChatRoom)
	router.PATCH("/prescriptions/:prescription_id", authMiddleware, userAuthorizationMiddleware, handler.SavePrescription)
	router.GET("/prescriptions", authMiddleware, userAuthorizationMiddleware, handler.GetAllPrescriptions)
	router.GET("/prescriptions/:prescription_id", authMiddleware, userAuthorizationMiddleware, handler.PreapereForCheckout)
	router.POST("/prescriptions/checkout", authMiddleware, userAuthorizationMiddleware, handler.CheckoutFromPrescription)
}

func corsRouting(router *gin.Engine, configCors cors.Config) {
	configCors.AllowAllOrigins = true
	configCors.AllowMethods = []string{"POST", "GET", "PUT", "PATCH", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept", "User-Agent", "Cache-Control"}
	configCors.ExposeHeaders = []string{"Content-Length"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))
}

func authenticationRouting(router *gin.Engine, handler *handler.AuthenticationHandler) {
	router.POST("/users/register", handler.RegisterUser)
	router.POST("/doctors/register", handler.RegisterDoctor)
	router.POST("/verification", handler.SendVerificationEmail)
	router.POST("/verification/password", handler.VerifyOneAccount)
	router.POST("/login", handler.Login)
	router.POST("/refresh-token", handler.GetNewAccessToken)
	router.POST("/reset-password", handler.SendResetPasswordToken)
	router.POST("/reset-password/verification", handler.ResetPasswordOneAccount)
}

func categoryRouting(router *gin.Engine, handler *handler.CategoryHandler, authMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc) {
	categoryRouter := router.Group("/categories")

	categoryRouter.GET("/", handler.GetAllCategories)
	categoryRouter.DELETE("/:category_id", authMiddleware, adminAuthorizationMiddleware, handler.DeleteCategory)
	categoryRouter.POST("/", authMiddleware, adminAuthorizationMiddleware, handler.AddOneCategory)
	categoryRouter.PUT("/:category_id", authMiddleware, adminAuthorizationMiddleware, handler.UpdateOneCategory)
}

func cartRouting(router *gin.Engine, handler *handler.CartHandler, authMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc) {
	cartRouter := router.Group("/cart")

	cartRouter.POST("/delivery", authMiddleware, userAuthorizationMiddleware, handler.CalculateDeliveryFee)
	cartRouter.POST("/", authMiddleware, userAuthorizationMiddleware, handler.CreateOneCart)
	cartRouter.PATCH("/:cart_id", authMiddleware, userAuthorizationMiddleware, handler.UpdateQtyCart)
	cartRouter.DELETE("/:cart_id", authMiddleware, userAuthorizationMiddleware, handler.DeleteOneCart)
	cartRouter.GET("/", authMiddleware, userAuthorizationMiddleware, handler.GetAllCart)
}

func orderRouting(router *gin.Engine, handler *handler.OrderHandler, authMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc, pharmacyManagerAuthorizationMiddleware gin.HandlerFunc) {
	router.POST("/orders", authMiddleware, userAuthorizationMiddleware, handler.CheckoutOrder)
	router.PATCH("/orders/:order_id/payment-proof", authMiddleware, userAuthorizationMiddleware, handler.UploadPaymentProofOrder)
	router.PATCH("/orders/:order_id/confirm-payment", authMiddleware, adminAuthorizationMiddleware, handler.ConfirmPayment)
	router.PATCH("/orders/:order_id/cancel-order", authMiddleware, userAuthorizationMiddleware, handler.CancelOrder)
	router.GET("/orders/:order_id", authMiddleware, userAuthorizationMiddleware, handler.GetOrderById)
	router.GET("/orders/pending", authMiddleware, userAuthorizationMiddleware, handler.GetAllUserPendingOrders)
	router.GET("/admin/orders", authMiddleware, adminAuthorizationMiddleware, handler.GetAllOrders)
}

func orderPharmacyRouting(router *gin.Engine, handler *handler.OrderPharmacyHandler, authMiddleware gin.HandlerFunc,
	pharmacyManagerAuthorizationMiddleware gin.HandlerFunc, userAuthorizationMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc) {
	router.PATCH("/pharmacy-orders/:order_pharmacy_id/send-package", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.UpdateStatusToSent)
	router.PATCH("/pharmacy-orders/:order_pharmacy_id/confirm-package", authMiddleware, userAuthorizationMiddleware, handler.UpdateStatusToConfirmed)
	router.PATCH("/pharmacy-orders/:order_pharmacy_id/cancel-package", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.UpdateStatusToCancelled)
	router.GET("/pharmacy-orders/:order_pharmacy_id", authMiddleware, userAuthorizationMiddleware, handler.GetOrderPharmacyById)
	router.GET("/pharmacy-orders", authMiddleware, userAuthorizationMiddleware, handler.GetAllUserOrderPharmacies)
	router.GET("/manager/pharmacy-orders", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetAllPartnerOrderPharmacies)
	router.GET("/manager/pharmacy-orders/summary", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetAllPartnerOrderPharmaciesSummary)
	router.GET("/admin/pharmacy-orders", authMiddleware, adminAuthorizationMiddleware, handler.GetAllOrderPharmacies)
}

func reportRouting(router *gin.Engine, handler *handler.ReportHandler, authMiddleware gin.HandlerFunc, pharmacyManagerAuthorizationMiddleware gin.HandlerFunc, adminAuthorizationMiddleware gin.HandlerFunc) {
	router.GET("/manager/categories/reports", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetPharmacyDrugCategoryReport)
	router.GET("/manager/drugs/reports", authMiddleware, pharmacyManagerAuthorizationMiddleware, handler.GetPharmacyDrugReport)
	router.GET("/admin/categories/reports", authMiddleware, adminAuthorizationMiddleware, handler.GetDrugCategoryReport)
	router.GET("/admin/drugs/reports", authMiddleware, adminAuthorizationMiddleware, handler.GetDrugReport)
}

func pingRouting(router *gin.Engine, handler *handler.PingHandler, authMiddleware gin.HandlerFunc, userMiddleware gin.HandlerFunc, doctorMiddleware gin.HandlerFunc, PharmacyManagerMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	pingRouter := router.Group("/ping")

	pingRouter.Use(authMiddleware)
	pingRouter.GET("all", handler.Ping)

	pingRouter.GET("/user", userMiddleware, handler.Ping)
	pingRouter.GET("/doctor", doctorMiddleware, handler.Ping)
	pingRouter.GET("/pharmacy-manager", PharmacyManagerMiddleware, handler.Ping)
	pingRouter.GET("/admin", adminMiddleware, handler.Ping)

	pingRouter.Use(userMiddleware)
	pingRouter.GET("/all-user-endpoints", handler.Ping)
}

func pprofRouting(router *gin.Engine) {
	pprofRouter := router.Group("/debug/pprof")

	pprofRouter.GET("/", gin.WrapH(http.HandlerFunc(pprof.Index)))
	pprofRouter.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
	pprofRouter.GET("/heap", gin.WrapH(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
	pprofRouter.GET("/block", gin.WrapH(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
	pprofRouter.GET("/goroutine", gin.WrapH(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
}
