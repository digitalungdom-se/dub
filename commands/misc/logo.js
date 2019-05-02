module.exports = {
  name: 'logo',
  description: 'Skickar serverns nuvarande logotyp',
  aliases: [ 'logotyp' ],
  group: 'misc',
  usage: 'logo',
  example: 'logo',
  serverOnly: false,
  adminOnly: false,
  execute( message, args ) {
    const logoURL = 'https://raw.githubusercontent.com/kelszo/dub/master/public/imgs/dub.png';

    message.reply( 'här kommer serverns och botens nuvarande logotyp', { file: logoURL } ).then( ( msg ) => msg.delete( 10000 ) );
  },
};