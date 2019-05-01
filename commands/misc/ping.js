module.exports = {
  name: 'ping',
  description: 'ping, pong!',
  aliases: [],
  group: 'misc',
  usage: 'ping',
  example: 'ping',
  serverOnly: false,
  adminOnly: false,
  execute( message, args ) {
    if ( !message.deleted ) message.delete( 10000 );
    const embed = {
      'title': ':ping_pong:',
      'description': `${new Date().getTime() - message.createdTimestamp}ms`,
      'color': 4086462
    };

    message.reply( { 'embed': embed } ).then( ( msg ) => msg.delete( 10000 ) );
  },
};