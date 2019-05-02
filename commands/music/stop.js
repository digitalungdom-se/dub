/* global guild controller client */

module.exports = {
  name: 'stop',
  description: 'Stannar botens nuvarande musik',
  aliases: [ 'stanna', 'st' ],
  group: 'music',
  usage: 'stop',
  example: 'stop',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( guild.voiceConnection ) {
      controller.stop();

      client.user.setActivity( 'Kelvin\'s cat', { type: 'WATCHING' } );
      return message.reply( 'stoppar boten' ).then( msg => { msg.delete( 10000 ); } );
    } else return message.reply( 'kan inte stanna boten då boten inte spelar något' ).then( msg => { msg.delete( 10000 ); } );
  },
};