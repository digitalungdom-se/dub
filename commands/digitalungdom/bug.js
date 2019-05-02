/* global db include guild */

const getUserByDiscordId = include( 'models/get' ).getUserByDiscordId;
const createNotificationEmbed = include( 'utils/embeds/createNotificationEmbed' );

module.exports = {
  name: 'bug',
  description: 'Skicka en bug report till Digital Ungdom.',
  aliases: [ 'bugg' ],
  group: 'digitalungdom',
  usage: 'bug <bugg>',
  example: 'bug verifierings funktionen funkar inte',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( args.length === 0 ) return message.reply( 'Du måste skicka med ett kort meddelande.' );
    const bug = args.join( ' ' );
    const authorID = message.author.id;
    const authorUsername = message.author.username;
    let id = await getUserByDiscordId( authorID );
    if ( id ) id = id._id;

    await db.collection( 'notifications' ).insertOne( {
      'type': 'bug',
      'where': 'discord',
      'message': bug,
      'author': {
        'id': id,
        'discordID': authorID,
        'discordUsername': authorUsername
      }
    } );

    const notification = createNotificationEmbed( 'BUG', bug, 16711680, { 'id': authorID, 'name': authorUsername, 'url': message.author.displayAvatarURL } );
    const notificationChannel = guild.channels.find( ch => ch.name === 'notifications' );
    notificationChannel.send( '@here', { 'embed': notification } );

    return message.reply( 'tack för din medverkan!' );
  },
};