/* global db include abs_path client cooldown */

const validateObjectID = require( 'mongodb' ).ObjectID.isValid;
const ObjectID = require( 'mongodb' ).ObjectID;
const fs = require( 'fs-extra' );
const util = require( 'util' );
const readFile = util.promisify( fs.readFile );
const Hogan = require( 'hogan.js' );
const checkDiscordID = include( 'models/check' ).checkDiscordID;

const sendMail = include( 'utils/sendMail' ).sendMail;

module.exports = {
  name: 'verify',
  description: 'Koppla ditt Discord-konto till ditt DU konto.',
  aliases: [ 'verifera' ],
  group: 'digitalungdom',
  usage: 'verify <username|code>',
  example: 'verify kelszo',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( message.channel.type === 'text' && !message.deleted ) message.delete();
    if ( args.length === 0 ) return message.author.send( 'Du måste skicka det användarnamn som du vill koppla till.' );
    const content = args[ 0 ];
    const discordID = message.author.id;

    const guild = await client.guilds.get( process.env.GUILD_ID );
    const user = await guild.fetchMember( message.author );
    if ( user.roles.find( r => r.name === 'verified' ) ) return message.author.send( 'Du är redan verifierad.' );
    if ( !( await checkDiscordID( message.author.id ) ).valid ) return message.author.send( 'Detta konto är redan kopplat till ett Digital Ungdom konto.' );

    if ( validateObjectID( content ) ) {
      const exists = ( await db.collection( 'users' ).findOneAndUpdate( { 'discordVerification': ObjectID( content ) }, {
        $set: { 'connectedApps.discord': discordID },
        $unset: { 'discordVerification': 1 },
      }, { projection: { '_id': 0, 'details.email': 1 } } ) ).value;
      if ( exists ) {
        const email = exists.details.email;

        const templateData = await readFile( abs_path( 'emails/verify/verifyConfirmation.mustache' ), 'utf8' );
        const template = Hogan.compile( templateData );
        const body = template.render( { name: message.author.username } );

        sendMail( email, 'Ditt Discord-konto är nu kopplat', body );

        const role = await guild.roles.find( r => r.name === 'verified' );
        const user = await guild.fetchMember( message.author );
        user.addRole( role );

        return message.author.send( 'Ditt konto är nu kopllat, grattis!' );
      } else message.author.send( 'Det finns inget konto kopplat till den id.' );
    } else {
      if ( new Date() < cooldown.verify[ message.author.id ] ) return message.author.send( 'Du får bara skicka ett verifikations mail per dag.' );
      const discordVerification = ObjectID();
      const search = { 'details.username': content.toLowerCase(), 'connectedApps.discord': { $exists: false } };
      const set = { $set: { 'discordVerification': discordVerification } };
      const projection = { projection: { '_id': 0, 'details.email': 1 }, collation: { locale: 'en', strength: 2 } };

      const exists = ( await db.collection( 'users' ).findOneAndUpdate( search, set, projection ) ).value;
      if ( exists ) {
        const email = exists.details.email;

        const templateData = await readFile( abs_path( 'emails/verify/verify.mustache' ), 'utf8' );
        const template = Hogan.compile( templateData );
        const body = template.render( { token: discordVerification } );

        let tomorrow = new Date();
        cooldown.verify[ message.author.id ] = tomorrow.setDate( tomorrow.getDate() + 1 );

        sendMail( email, 'Koppla ditt Discord-konto', body );
        return message.author.send( 'Ett email har skickats till kontot kopplat till det användarnamnet!' );
      } else {
        return message.author.send( 'Det gick inte att hitta ett konto, detta kan bero på tre saker:\n**1.** Det finns inget konto med det användarnamn.\n**2.** Du har redan skickat ett email, testa igen om 24 timmar.\n**3.** Du är redan verifierad.' );
      }
    }
  },
};