/* global controller guild */

module.exports = {
  name: 'skip',
  description: 'Skippar den nuvarande låt',
  aliases: [ 'skippa', 'byt', 'sk' ],
  group: 'music',
  usage: 'skip',
  example: 'skip',
  serverOnly: true,
  adminOnly: false,
  async execute( message, args ) {
    if ( guild.voiceConnection && controller.queue.length > 0 ) {
      controller.skip();
      if ( controller.queue.length === 0 ) return message.reply( 'Kön är slut.' ).then( msg => { msg.delete( 10000 ); } );
      else {
        message.reply( 'skippar låten' ).then( msg => { msg.delete( 10000 ); } );
      }
    } else return message.reply( 'kan inte skippa låt då boten inte spelar något.' ).then( msg => { msg.delete( 10000 ); } );
  },
};