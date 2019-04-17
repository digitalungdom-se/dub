module.exports = {
  name: 'member',
  description: 'Skickar bli medlem länken.',
  aliases: [ 'medlem', 'bli-medlem' ],
  group: 'digitalungdom',
  usage: 'member',
  serverOnly: false,
  execute( message, args ) {
    message.reply( 'här är länken att bli medlem: https://digitalungdom.se/bli-medlem' );
  },
};