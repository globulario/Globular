package ressource_client

import (
	//"encoding/json"
	"log"
	"testing"
)

var (
	// Connect to the plc client.
	client = NewRessource_Client("10.67.44.131", "ressource.RessourceService")
)

/*
// Test various function here.
func TestRegisterAccount(t *testing.T) {

	log.Println("---> test register a new account.")
	err := client.RegisterAccount("davecourtois", "dave.courtois60@gmail.com", "1234", "1234")
	if err != nil {
		log.Println("---> ", err)
	}
}
*/
func TestAuthenticate(t *testing.T) {
	log.Println("---> test authenticate account.")
	//token, err := client.Authenticate("dave.courtois60@gmail.com", "1234")
	token, err := client.Authenticate("sa", "adminadmin")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", token)
	}
}

/*
func TestCreateRole(t *testing.T) {
	log.Println("---> create role ")
	err := client.CreateRole("db_user", []string{"/persistence.PersistenceService/InsertOne", "/persistence.PersistenceService/InsertMany", "/persistence.PersistenceService/Find", "/persistence.PersistenceService/FindOne", "/persistence.PersistenceService/ReplaceOne", "/persistence.PersistenceService/DeleteOne", "/persistence.PersistenceService/Delete", "/persistence.PersistenceService/Count", "/persistence.PersistenceService/Update", "/persistence.PersistenceService/UpdateOne"})
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestAddRoleAction(t *testing.T) {
	log.Println("---> Add Role action ")
	err := client.AddRoleAction("db_user", "/toto")
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestRemoveRoleAction(t *testing.T) {
	log.Println("---> Remove Role action ")
	err := client.RemoveRoleAction("db_user", "/toto")
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestAddAccountRole(t *testing.T) {
	log.Println("---> Add account Role ")
	err := client.AddAccountRole("davecourtois", "db_user")
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestRemoveAccountRole(t *testing.T) {
	log.Println("---> Remove account Role ")
	err := client.RemoveAccountRole("davecourtois", "db_user")
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestDeleteRole(t *testing.T) {
	log.Println("---> Delete role ")
	err := client.DeleteRole("db_user")
	if err != nil {
		log.Println("---> ", err)
	}
}

// Remove an account.
func TestDeleteAccount(t *testing.T) {

	log.Println("---> test remove existing account.")
	err := client.DeleteAccount("davecourtois")
	if err != nil {

		log.Println("---> ", err)
	}
}
*/

/** Set file permission **/
/*func TestSetPermission(t *testing.T) {
	log.Println("---> Set permission by user")
	err := client.SetFilePermissionByUser("davecourtois", "test", 6)
	if err != nil {

		log.Println("---> ", err)
	}

	err = client.SetFilePermissionByRole("guest", "test/sub-folder/another-text.txt", 6)
	if err != nil {

		log.Println("---> ", err)
	}
}*/

/** Display the created file permission **/
/*func TestGetPermissions(t *testing.T) {
	log.Println("---> Get permissions by user")
	permissions, err := client.GetFilePermissions("test")
	if err != nil {
		log.Println("---> ", err)
		return
	}
	for i := 0; i < len(permissions); i++ {
		log.Println(permissions[i])
	}

	permissions, err = client.GetFilePermissions("test/sub-folder/another-text.txt")
	if err != nil {
		log.Println("---> ", err)
		return
	}
	for i := 0; i < len(permissions); i++ {
		log.Println(permissions[i])
	}
}*/

/** Display the created file permission **/
/*func TestDeletePermissions(t *testing.T) {
	log.Println("---> Delete permissions")
	err := client.DeleteFilePermissions("test", "")
	if err != nil {
		log.Println("---> ", err)
		return
	}
}*/

/** Get the root file informations **/
/*func TestGetAllFilesInfo(t *testing.T) {
	log.Println("---> Get All File Info")
	infos, err := client.GetAllFilesInfo()
	if err != nil {
		log.Println("---> ", err)
		return
	}
	log.Println(infos)
}*/
