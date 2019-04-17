module.exports = {
  name: 'logo',
  description: 'Skickar den nuvarande server logotyp',
  aliases: [ 'logotyp' ],
  group: 'misc',
  usage: 'logo',
  serverOnly: false,
  execute( message, args ) {
    const logoURL = 'https://raw.githubusercontent.com/kelszo/dub/master/dub.png';

    message.reply( 'här kommer serverns nuvarande logotyp', { file: logoURL } );
  },
};