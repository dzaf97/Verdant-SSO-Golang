package network

type (
	// CREDENTIAL MANAGEMENT
	ForgotPasswordReq struct {
		Email string
	}

	ResetPasswordReq struct {
		Token    string
		Password string
	}

	// USER MANAGEMENT
	ListUser struct {
		Username       string
		RoleName       string
		DateRegistered string
	}

	UserDetail struct {
		FirstName     string
		LastName      string
		RoleID        int
		PhoneNo       string
		Email         string
		Address       string
		StateName     string
		Postcode      string
		HouseholdSize int
		BuildupArea   int
		NumberOfRooms int
		PanelBrand    string
		Dsy           int
	}

	AddRole struct {
		RoleName string
	}

	GetRole struct {
		RoleID   int
		RoleName string
	}
)
