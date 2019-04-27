/* global config */

module.exports = {
  name: 'clear',
  description: 'Rensar alla meddelande som har med botten att göra.',
  aliases: [ 'rensa' ],
  group: 'admin',
  usage: 'clear',
  example: 'clear',
  serverOnly: true,
  adminOnly: true,
  async execute( message, args ) {
    message.delete();
    let messages = await message.channel.fetchMessages( { limit: 100 } );
    messages = messages.filter( function ( m ) {
      if ( m.author.bot ) return true;
      else if ( config.prefix.indexOf( m.content.charAt( 0 ) ) !== -1 ) return true;
    } );
    message.channel.bulkDelete( messages );
  },
};