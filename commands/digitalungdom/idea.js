/* global db include guild */

const getUserByDiscordId = include( 'models/get' ).getUserByDiscordId;
const createNotificationEmbed = include( 'utils/embeds/createNotificationEmbed' );

module.exports = {
  name: 'idea',
  description: 'Skicka in ett förslag till Digital Ungdom.',
  aliases: [ 'suggestion', 'förslag' ],
  group: 'digitalungdom',
  usage: 'idea <idea>',
  example: 'idea kommentera er kod bättre',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( args.length === 0 ) return message.reply( 'Du måste skicka med ett kort meddelande.' );
    const idea = args.join( ' ' );
    const authorID = message.author.id;
    const authorUsername = message.author.username;
    let id = await getUserByDiscordId( authorID );
    if ( id ) id = id._id;

    await db.collection( 'notifications' ).insertOne( {
      'type': 'idea',
      'where': 'discord',
      'message': idea,
      'author': {
        'id': id,
        'discordID': authorID,
        'discordUsername': authorUsername
      }
    } );

    const notification = createNotificationEmbed( 'IDEA', idea, 65397, { 'id': authorID, 'name': authorUsername, 'url': message.author.displayAvatarURL } );
    const notificationChannel = guild.channels.find( ch => ch.name === 'notifications' );
    notificationChannel.send( '@here', { 'embed': notification } );

    return message.reply( 'tack för din medverkan!' ).then( msg => { msg.delete( 10000 ); } );
  },
};