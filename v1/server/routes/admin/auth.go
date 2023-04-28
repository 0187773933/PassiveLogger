package adminroutes

import (
	"fmt"
	"time"
	"strings"
	// "reflect"
	// "unsafe"
	fiber "github.com/gofiber/fiber/v2"
	bcrypt "golang.org/x/crypto/bcrypt"
	encryption "github.com/0187773933/PassiveLogger/v1/encryption"
)

func validate_login_credentials( context *fiber.Ctx ) ( result bool ) {
	result = false
	uploaded_username := context.FormValue( "username" )
	if uploaded_username == "" { fmt.Println( "username empty" ); return }
	if uploaded_username != GlobalConfig.AdminUsername { fmt.Println( "username not correct" ); return }
	uploaded_password := context.FormValue( "password" )
	if uploaded_password == "" { fmt.Println( "password empty" ); return }
	password_matches := bcrypt.CompareHashAndPassword( []byte( uploaded_password ) , []byte( GlobalConfig.AdminPassword ) )
	if password_matches != nil { fmt.Println( "bcrypted password doesn't match" ); return }
	result = true
	return
}

func Logout( context *fiber.Ctx ) ( error ) {
	context.Cookie( &fiber.Cookie{
		Name: GlobalConfig.ServerCookieName ,
		Value: "" ,
		Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
		HTTPOnly: true ,
		Secure: true ,
	})
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Out</h1>" )
}

// r.Header.collectCookies()
// r.Header.DelAllCookies()
func clear_cookies( context *fiber.Ctx ) {
	cookie_header := context.Request().Header.Peek( "Cookie" )
	raw_cookies := strings.Split( string( cookie_header ) , ";" )
	for _ , raw_cookie := range raw_cookies {
		parts := strings.SplitN( raw_cookie , "=" , 2 )
		if len( parts ) < 1 { continue }
		context.Cookie( &fiber.Cookie{
			Name: parts[ 0 ] ,
			Value: "" ,
			Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
			HTTPOnly: true ,
			Secure: true ,
		})
	}
}

// POST http://localhost:5950/admin/login
func HandleLogin( context *fiber.Ctx ) ( error ) {
	valid_login := validate_login_credentials( context )
	if valid_login == false { return serve_failed_attempt( context ) }
	clear_cookies( context )
	cookie := fiber.Cookie{
		Name: GlobalConfig.ServerCookieName ,
		Value: encryption.SecretBoxEncrypt( GlobalConfig.BoltDBEncryptionKey , GlobalConfig.ServerCookieSecretMessage ) ,
		Secure: true ,
		Path: "/" ,
		// Domain: "blah.ngrok.io" , // probably should set this for webkit
		HTTPOnly: true ,
		SameSite: "Lax" ,
		Expires: time.Now().AddDate( 10 , 0 , 0 ) , // aka 10 years from now
	}
	fmt.Println( "Valid Login , Setting Cookie" , cookie )
	context.Cookie( &cookie )
	return context.Redirect( "/" )
}

func validate_admin_cookie( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( GlobalConfig.ServerCookieName )
	if admin_cookie == "" { fmt.Println( "admin cookie was blank" ); return }
	admin_cookie_value := encryption.SecretBoxDecrypt( GlobalConfig.BoltDBEncryptionKey , admin_cookie )
	if admin_cookie_value != GlobalConfig.ServerCookieSecretMessage { fmt.Println( "cookie secret message was not equal" ); return }
	result = true
	return
}

func serve_failed_attempt( context *fiber.Ctx ) ( error ) {
	// context.Set( "Content-Type" , "text/html" )
	// return context.SendString( "<h1>no</h1>" )
	return context.SendFile( "./v1/server/html/login.html" )
}

func ServeLoginPage( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/login.html" )
}

func ServeAuthenticatedPage( context *fiber.Ctx ) ( error ) {
	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
	x_path := context.Route().Path
	url_key := strings.Split( x_path , "/admin" )
	if len( url_key ) < 2 { return context.SendFile( "./v1/server/html/login.html" ) }
	// fmt.Println( "Sending -->" , url_key[ 1 ] , x_path )
	return context.SendFile( ui_html_pages[ url_key[ 1 ] ] )
}

