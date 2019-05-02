/* global version live lastUpdated */

module.exports = {
  name: 'info',
  description: 'Anger information om boten.',
  group: 'misc',
  usage: 'info',
  example: 'info',
  serverOnly: false,
  adminOnly: false,
  execute( message, args ) {
    if ( !message.deleted ) message.delete( 10000 );
    const embed = {
      'description': '__**INFORMATION OM BOTEN**__',
      'color': 4086462,
      'fields': [ {
          'name': 'VERSION',
          'value': `${version} (${lastUpdated})`,
          'inline': true
        },
        {
          'name': 'LIVE SEDAN',
          'value': `${live}`,
          'inline': true
        },
        {
          'name': 'KÄLLKOD',
          'value': '[github](https://github.com/kelszo/dub)',
          'inline': true
        },
        {
          'name': 'MEDARBETARE',
          'value': '<@217632464531619852>, <@228889878861971456>',
          'inline': true
        }
      ]
    };

    message.reply( { embed: embed } );
  },
};