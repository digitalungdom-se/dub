/* global db include abs_path client guild */

const validateObjectID = require( 'mongodb' ).ObjectID.isValid;
const ObjectID = require( 'mongodb' ).ObjectID;
const fs = require( 'fs-extra' );
const util = require( 'util' );
const readFile = util.promisify( fs.readFile );
const Hogan = require( 'hogan.js' );

const sendMail = include( 'utils/sendMail' ).sendMail;

module.exports = {
  name: 'verify',
  description: 'Koppla ditt discord konto till ditt Digital Ungdom konto.',
  aliases: [ 'verifera' ],
  group: 'digitalungdom',
  usage: 'verify <username|code>',
  example: 'verify kelszo',
  serverOnly: false,
  async execute( message, args ) {
    if ( message.channel.type === 'text' ) message.delete();
    if ( args.length === 0 ) return message.author.send( 'Du måste skicka det användarnamn som du vill koppla till.' );
    const content = args[ 0 ];
    const discordId = message.author.id;

    const guild = await client.guilds.get( process.env.GUILD_ID );
    const user = await guild.fetchMember( message.author );
    if ( user.roles.find( r => r.name === 'verified' ) ) return message.author.send( 'Du är redan verifierad.' );

    if ( validateObjectID( content ) ) {
      const exists = ( await db.collection( 'users' ).findOneAndUpdate( { 'discordVerification': content }, { $set: { 'discordId': discordId }, $unset: { 'discordVerification': 1, 'cooldowns.dub.verify': 1 } }, { projection: { '_id': 0, 'email': 1 } } ) ).value;
      if ( exists ) {
        const email = exists.email;

        const templateData = await readFile( abs_path( 'emails/verify/verifyConfirmation.mustache' ), 'utf8' );
        const template = Hogan.compile( templateData );
        const body = template.render( { name: message.author.username } );

        await sendMail( email, 'Ditt Discord konto är nu kopplat', body );

        const role = await guild.roles.find( r => r.name === 'verified' );
        const user = await guild.fetchMember( message.author );
        await user.addRole( role );

        return message.author.send( 'Ditt konto är nu kopllat, grattis!' );
      } else message.author.send( 'Det finns inget konto kopplat till den id.' );
    } else {
      const discordVerification = ObjectID().toString();
      const search1 = { 'usernameLower': content.toLowerCase(), 'discordId': { $exists: false }, 'cooldowns.dub.verify': { $gte: new Date() } };
      const search2 = search1;
      search2[ 'cooldowns.dub.verify' ] = { $exists: false };

      const exists = ( await db.collection( 'users' ).findOneAndUpdate( { $or: [ search1, search2 ] }, { $set: { 'discordVerification': discordVerification, 'cooldowns.dub.verify': new Date( ( new Date ).getTime() + 86400000 ) } }, { projection: { '_id': 0, 'email': 1 } } ) ).value;
      if ( exists ) {
        const email = exists.email;

        const templateData = await readFile( abs_path( 'emails/verify/verify.mustache' ), 'utf8' );
        const template = Hogan.compile( templateData );
        const body = template.render( { token: discordVerification } );

        await sendMail( email, 'Koppla ditt Discord konto', body );
        return message.author.send( 'Ett email har skickats till kontot kopplat till det användarnamnet!' );
      } else {
        return message.author.send( 'Det gick inte att hitta ett konto, detta kan bero på tre saker:\n**1.** Det finns inget konto med det användarnamn.\n**2.** Du har redan skickat ett email, testa igen om 24 timmar.\n **3.** Du är redan verifierad.' );
      }
    }
  },
};