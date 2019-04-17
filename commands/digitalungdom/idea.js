/* global db */

module.exports = {
  name: 'idea',
  description: 'Skicka ett förslag till Digital Ungdom. Det kan vara t.ex. ett förbättrings förslag eller önskemål.',
  aliases: [ 'suggestions', 'förslag' ],
  group: 'digitalungdom',
  usage: 'idea <idea>',
  example: 'idea skriva kommentarer i er kod',
  serverOnly: false,
  async execute( message, args ) {
    if ( args.length === 0 ) return message.reply( 'Du måste skicka med ett kort meddelande.' );
    const idea = args.join( ' ' );
    const authorId = message.author.id;
    const authorUsername = message.author.username;
    let id = await db.collection( 'users' ).findOne( { 'discordId': authorId }, { projection: { '_id': 1 } } );
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