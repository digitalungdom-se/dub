module.exports = {
  name: 'slap',
  description: 'Smiskar en styrelsemedlem.',
  aliases: [ 'smisk' ],
  group: 'misc',
  usage: 'slap <medlem>',
  example: 'slap simon',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( args.length === 0 ) return message.reply( 'Du måste ge mig ett namn för att smiska.' ).then( ( msg ) => msg.delete( 10000 ) );
    const name = args[ 0 ].toLowerCase();

    const members = {
      'kelvin': '217632464531619852',
      'simon': '228889878861971456',
      'douglas': '297671552823066626',
    };
    if ( !members[ name ] ) return message.reply( 'ingen i styrelsen heter så, testa: Douglas, Simon eller Kelvin' ).then( ( msg ) => msg.delete( 10000 ) );
    message.channel.send( `Du har varit en riktig stygg pojk <@${members[name]}>. *smisk*` );
  },
};