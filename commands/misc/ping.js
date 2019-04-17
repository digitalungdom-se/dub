module.exports = {
  name: 'ping',
  description: 'ping, pong!',
  aliases: [],
  group: 'misc',
  usage: 'ping',
  serverOnly: false,
  execute( message, args ) {
    message.reply( 'pong.' );
  },
};