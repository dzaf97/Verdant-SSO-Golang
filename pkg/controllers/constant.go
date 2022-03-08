package controllers

import "time"

var loc, _ = time.LoadLocation("Asia/Kuala_Lumpur")

//should be banyak or spaerate , idk
//think on error logging streaming
const (
	ErrInvalidDataStruct = "Unsupported data format"
	ErrProcessError      = "Something not right"
	ErrDbTransFail       = "Error during transaction"
	ErrDbInsertFail      = "Data transaction failed"
	ProcessComplete      = "Data transaction complete"
	// ErrReadJSON          = "Cannot read data format"

	ErrValidate         = "Data validation failed: Invalid Parameters."
	ErrRoleNameNotAvail = "Role name not available"
	ErrStatNameNotAvai  = "Department name not available"

	//Login
	ErrAuthFail      = "Wrong username or password"
	ErrAuthSuspended = "Account suspended. Please contact admin"
	ErrEmailUsed     = "Email not available"
	ProcinCmplt      = "Data successfully insert"
	ProcinUpdt       = "Data successfully updated"
	RecNotFound      = "Record not found"
	RecIsFound       = "Record available"
	ErrInvalidParam  = "Invalid parameter"
	FailUpdate       = "Information can be update"
	RegSuccess       = "Staff registration complate"

	//acl
	DntHvePmr = "You dont have permission to access this resources"
	AclAssco  = "There is accout asscited with this ACL"
	//http
	MthdNotAllw = "Method Not Allowed"

	//logout
	SuccessLogOut = "Logout success"

	//Override
	OverOK   = "Scheduler Override success"
	OverAvai = "Scheduler override is ongoing"
	OverFail = "Scheduler override failed"

	//attr
	EtyNotFnd   = "Entity not found"
	TypeDevice  = "DEVICE"
	TypeCluster = "CLUSTER"
	TypeTenant  = "TENANT"
	TypeUser    = "USER"
)
