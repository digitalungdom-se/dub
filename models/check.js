/* global db */

const validator = require( 'validator' );

module.exports.checkUsername = async function ( username ) {
  const userExists = ( await db.collection( 'users' ).find( { 'details.username': username.toLowerCase() }, { 'projection': { '_id': 1 } } ).collation( { locale: 'en', strength: 2 } ).toArray() )[ 0 ];

  if ( userExists ) return { 'valid': false, 'field': 'username' };
  else return { 'valid': true, 'field': 'username' };
};

module.exports.checkEmail = async function ( email ) {
  // Normalises the email using validatorjs. See documentation for exact normalisation rules
  email = validator.normalizeEmail( email );

  const userExists = await db.collection( 'users' ).findOne( { 'details.email': email }, { 'projection': { '_id': 1 } } );

  if ( userExists ) return { 'valid': false, 'field': 'email' };
  else return { 'valid': true, 'field': 'email' };
};

module.exports.checkDiscordID = async function ( id ) {
  const userExists = await db.collection( 'users' ).findOne( { 'connectedApps.discord': id }, { 'projection': { '_id': 1 } } );

  if ( userExists ) return { 'valid': false, 'field': 'email' };
  else return { 'valid': true, 'field': 'email' };
};