module.exports = {
  name: 'ping',
  description: 'ping, pong!',
  aliases: [],
  group: 'misc',
  usage: 'ping',
  execute( message, args ) {
    message.reply( 'Pong.' );
  },
};