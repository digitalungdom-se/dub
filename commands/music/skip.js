/* global player musicQueue */

module.exports = {
  name: 'skip',
  description: 'Skippar den nuvarande låt',
  aliases: [ 'skippa', 'byt', 'sk' ],
  group: 'music',
  usage: 'skip',
  serverOnly: true,
  execute( message, args ) {
    if ( message.guild.voiceConnection ) {
      if ( musicQueue.length === 0 ) message.channel.send( 'Kön är slut.' );
      else message.reply( 'skippar låten' );
      player.end();
    } else message.reply( 'kan inte skippa låt då botten inte spelar något.' );
  },
};