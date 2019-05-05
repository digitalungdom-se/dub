/* global config */

module.exports = {
  name: 'clear',
  description: 'Rensar alla meddelande till och från boten.',
  aliases: [ 'rensa' ],
  group: 'admin',
  usage: 'clear <all>',
  example: 'clear',
  serverOnly: true,
  adminOnly: true,
  async execute( message, args ) {
    let messages = await message.channel.fetchMessages( { limit: 100 } );
    if ( args[ 0 ] !== 'all' ) {
      messages = messages.filter( function ( m ) {
        if ( m.author.bot ) return true;
        else if ( config.prefix.indexOf( m.content.charAt( 0 ) ) !== -1 ) return true;
      } );
    }
    message.channel.bulkDelete( messages );

    if ( message.channel.name === 'voting' ) global.voteDic = {};
  },
};