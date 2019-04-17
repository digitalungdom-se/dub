/* global player musicQueue */

module.exports = {
  name: 'pr',
  description: 'Pausar eller återupptar musiken.',
  aliases: [ 'pause', 'resume', 'pausa', 'fortsätt' ],
  group: 'music',
  usage: 'pr',
  serverOnly: true,
  execute( message, args ) {
    if ( !message.guild.voiceConnection || !player || musicQueue.length === 0 ) return message.reply( 'botten spelar inget.' );
    if ( player.paused ) player.resume();
    else player.pause();
  },
};