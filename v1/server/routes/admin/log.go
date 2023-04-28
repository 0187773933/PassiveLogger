package adminroutes

import (
	"fmt"
	"time"
	binary "encoding/binary"
	// json "encoding/json"
	net_url "net/url"
	fiber "github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	bolt "github.com/boltdb/bolt"
	encryption "github.com/0187773933/PassiveLogger/v1/encryption"
	utils "github.com/0187773933/PassiveLogger/v1/utils"
)

func get_context_uuid( context *fiber.Ctx ) ( result string ) {
	uploaded_uuid := context.Params( "uuid" )
	parsed_uuid , err := uuid.FromString( uploaded_uuid )
	if err != nil { return }
	if parsed_uuid.Version() != uuid.V4 { return }
	result = parsed_uuid.String()
	return
}

func return_fail( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Failed To Log</h1>" )
}

// https://github.com/boltdb/bolt#autoincrementing-integer-for-the-bucket
// itob returns an 8-byte big endian representation of v.
func itob( v uint64 ) []byte {
    b := make( []byte , 8 )
    binary.BigEndian.PutUint64( b , uint64( v ) )
    return b
}

func LogMessage( context *fiber.Ctx ) ( error ) {

	// Validation
	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
	context_uuid := get_context_uuid( context )
	if context_uuid == "" { fmt.Println( "no valid uuid was sent in url" ); return return_fail( context ) }
	uploaded_message := context.Params( "message" )
	if uploaded_message == "" { return return_fail( context ) }

	unescaped_message , _ := net_url.QueryUnescape( uploaded_message )
	time_string := utils.GetFormattedTimeString()
	log_message := fmt.Sprintf( "%s === %s" , time_string , unescaped_message )
	fmt.Println( "Logging:" , log_message )
	encrypted_log_message := encryption.ChaChaEncryptBytes( GlobalConfig.BoltDBEncryptionKey , []byte( log_message ) )

	db , _ := bolt.Open( GlobalConfig.BoltDBPath , 0600 , &bolt.Options{ Timeout: ( 3 * time.Second ) } )
	defer db.Close()
	db_result := db.Update( func( tx *bolt.Tx ) error {
		uuid_bucket , _ := tx.CreateBucketIfNotExists( []byte( context_uuid ) )
		sequence_id  , _ := uuid_bucket.NextSequence()
		fmt.Println( "sequence id" , sequence_id )
		uuid_bucket.Put( itob( sequence_id ) , encrypted_log_message )
		return nil
	})
	if db_result != nil { return return_fail( context ) }
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Whatever it Was</h1>" )
}

func LogViewMessages( context *fiber.Ctx ) ( error ) {
	// Validation
	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
	context_uuid := get_context_uuid( context )
	if context_uuid == "" { fmt.Println( "no valid uuid was sent in url" ); return return_fail( context ) }

	db , _ := bolt.Open( GlobalConfig.BoltDBPath , 0600 , &bolt.Options{ Timeout: ( 3 * time.Second ) } )
	defer db.Close()

	var messages [][]byte
	db_result := db.View( func( tx *bolt.Tx ) error {
		uuid_bucket := tx.Bucket( []byte( context_uuid ) )
		uuid_bucket.ForEach( func( sequence_id , encrypted_message []byte ) error {
			decrytped_message := encryption.ChaChaDecryptBytes( GlobalConfig.BoltDBEncryptionKey , encrypted_message )
			messages = append( messages , decrytped_message )
			return nil
		})
		return nil
	})
	if db_result != nil { return return_fail( context ) }

	return context.JSON( fiber.Map{
		"route": fmt.Sprintf( "/v/%s" , context_uuid ) ,
		"messages": messages ,
	})
}

// func LogObject( context *fiber.Ctx ) ( error ) {
// 	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
// 	context_uuid := get_context_uuid( context )
// 	if context_uuid == "" { fmt.Println( "no valid uuid was sent in url" ); return return_fail( context ) }
// 	context.Set( "Content-Type" , "text/html" )
// 	return context.SendString( "<h1>Logged Whatever it was</h1>" )
// }

// uuid "github.com/satori/go.uuid"
// uuid.NewV4().String()
// uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440000")