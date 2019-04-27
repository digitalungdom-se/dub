/* global db include */

const getUserByDiscordId = include( 'models/get' ).getUserByDiscordId;

module.exports = {
  name: 'idea',
  description: 'Skicka ett förslag till Digital Ungdom. Det kan vara t.ex. ett förbättrings förslag eller önskemål.',
  aliases: [ 'suggestion', 'förslag', 'id' ],
  group: 'digitalungdom',
  usage: 'idea <idea>',
  example: 'idea skriva kommentarer i er kod',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( args.length === 0 ) return message.reply( 'Du måste skicka med ett kort meddelande.' );
    const idea = args.join( ' ' );
    const authorId = message.author.id;
    const authorUsername = message.author.username;
    let id = await getUserByDiscordId( authorId );
    if ( id ) id = id._id;

    await db.collection( 'notifications' ).insertOne( {
      'type': 'idea',
      'message': idea,
      'author': {
        'id': id,
        'discordId': authorId,
        'discordUsername': authorUsername
      }
    } );

    return message.reply( 'tack för din medverkan!' );
  },
};