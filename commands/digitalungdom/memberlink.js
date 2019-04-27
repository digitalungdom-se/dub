module.exports = {
  name: 'memberlink',
  description: 'Skickar bli medlem länken.',
  aliases: [ 'medlemlänk', 'ml' ],
  group: 'digitalungdom',
  usage: 'member',
  serverOnly: false,
  adminOnly: false,
  execute( message, args ) {
    message.reply( 'här är länken att bli medlem: https://digitalungdom.se/bli-medlem' );
  },
};