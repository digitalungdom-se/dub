/* global controller guild */

module.exports = {
  name: 'pr',
  description: 'Pausar eller återupptar musiken.',
  aliases: [ 'pause', 'resume', 'pausa', 'fortsätt' ],
  group: 'music',
  usage: 'pr',
  example: 'pr',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( !guild.voiceConnection || !controller || controller.queue.length === 0 ) return message.reply( 'botten spelar inget.' ).then( msg => { msg.delete( 10000 ); } );
    controller.pauseResume();
  },
};