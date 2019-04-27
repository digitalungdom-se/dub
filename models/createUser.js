/* global db abs_path include */

const bcrypt = require( 'bcryptjs' );
const crypto = require( 'crypto' );
const fs = require( 'fs-extra' );
const util = require( 'util' );
const readFile = util.promisify( fs.readFile );
const Hogan = require( 'hogan.js' );

const sendMail = include( 'utils/sendMail' ).sendMail;

module.exports = async function createUser( user ) {
  // Hahes the users password with 12 salt rounds (standard)
  user.details.password = await bcrypt.hash( user.details.password, 13 );

  // Generates 32 character (byte) long token to use for as a verification token. Inserts it into the users mongodb document and send them an email.
  const verificationToken = crypto.randomBytes( 32 ).toString( 'hex' );
  user.verificationToken = verificationToken;

  const templateData = await readFile( abs_path( 'emails/register/verifyEmail.mustache' ), 'utf8' );
  const template = Hogan.compile( templateData );
  const body = template.render( { token: verificationToken } );

  await Promise.all( [
    sendMail( user.details.email, 'Verifiera din e-postadress', body ),
    db.collection( 'users' ).insertOne( user ),
  ] );

  return { 'error': false };
};