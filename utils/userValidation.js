/* global include */

const validator = require( 'validator' );

const checkUsername = include( 'models/check' ).checkUsername;
const checkEmail = include( 'models/check' ).checkEmail;

module.exports = {
  'details': {
    'username': async function ( username ) {
      if ( typeof username != 'string' ) return { 'error': 'användarnamn är inte en sträng', 'return': username, 'field': 'username' };
      // Validates username according to following rules: min 3 max 24 characters and only includes valid characters (A-Z, a-z, 0-9, and _)
      if ( !validator.isLength( username, { min: 3, max: 24 } ) ) return { 'error': 'användarnamn är inte inom ramen 3-24 bokstäver', 'return': username, 'field': 'username' };
      if ( !/^(\w+)$/.test( username ) ) return { 'error': 'icketillåtna bokstäver i användarnamnet', 'return': username, 'field': 'username' };

      const usernameExists = await checkUsername( username );
      if ( !usernameExists.valid ) return { 'error': 'användarnamnet finns redan', 'return': username, 'field': 'username' };

      return { 'error': false };
    }, 'email': async function ( email ) {
      if ( typeof email != 'string' ) return { 'error': 'e-posten är inte en sträng', 'return': email, 'field': 'email' };
      // Validates email according to following rules: is a valid email.
      if ( !validator.isEmail( email ) ) return { 'error': 'felaktig e-postadress', 'return': email, 'field': 'email' };
      // Normalises email according to validatorjs (see validatorjs documentation for rules)
      email = validator.normalizeEmail( email );

      const usernameExists = await checkEmail( email );
      if ( !usernameExists.valid ) return { 'error': 'e-posten finns redan', 'return': email, 'field': 'email' };

      return { 'error': false };
    }, 'password': async function ( password ) {
      if ( typeof password != 'string' ) return { 'error': 'lösenordet är inte en sträng', 'return': password, 'field': 'password' };
      // Validates password according to following rules: min 8 max 72 characters, includes at least one character and one number
      if ( !validator.isLength( password, { min: 8, max: 72 } ) ) return { 'error': 'lösenordet är inte inom ramen 8-72 karaktärer', 'return': password, 'field': 'password' };
      if ( !/((.*[a-öA-Ö])(.*[0-9]))|((.*[0-9])(.*[a-öA-Ö]))/.test( password ) ) return { 'error': 'lösenordet är inte tillräkligt stark', 'return': password, 'field': 'password' };

      return { 'error': false };
    }, 'name': async function ( name ) {
      if ( typeof name != 'string' ) return { 'error': 'namn är inte en sträng', 'return': name, 'field': 'name' };

      // Validates name according to following rules: min 3 max 64 characters, min 2 names (e.g. Firstname Surname), only includes allowed characters (A-Z, a-z (including all diatrics), and - ' , . ')
      if ( !validator.isLength( name, { min: 3, max: 64 } ) ) return { 'error': 'namn är inte inom ramen 3-24 bokstäver', 'return': name, 'field': 'name' };
      if ( name.split( ' ' ).filter( n => n ).length < 2 ) return { 'error': 'inte tillräkligt många namn, måste innehålla minst för- och efternamn', 'return': name, 'field': 'name' };
      if ( !/^(([A-Za-zÀ-ÖØ-öø-ÿ\-',.\s ]+))$/.test( name ) ) return { 'error': 'icketillåtna bokstäver i namnet', 'return': name, 'field': 'name' };

      return { 'error': false };
    }, 'birthdate': async function ( birthdate ) {
      if ( typeof birthdate != 'string' ) return { 'error': 'födelsedag är inte en sträng', 'return': birthdate, 'field': 'birthdate' };

      // Validates birthdate according to following rules: makes sure that the date is correct length, makes sure that is is a date (strict, i.e. that is is a valid date too. See validatorjs documentation), and that is actually is a birthdate (i.e. is before the current date).
      if ( !validator.isISO8601( birthdate, { strict: true } ) ) return { 'error': 'felaktig födelsedag', 'return': birthdate, 'field': 'birthdate' };
      if ( !validator.isBefore( birthdate ) ) return { 'error': 'back to the future?', 'return': birthdate, 'field': 'birthdate' };

      return { 'error': false };
    },
  },
};