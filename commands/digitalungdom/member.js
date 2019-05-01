/* global memberProcess include */
const validate = include( 'utils/userValidation' );
const checkDiscordId = include( 'models/check' ).checkDiscordId;
const createUser = include( 'models/createUser' );

const user = {
  'details': {
    'username': '',
    'email': '',
    'verified': false,
    'password': '',
    'name': '',
    'birthdate': '',
    'gender': '',
    'profilePicture': null,
    'school': undefined,
    'company': undefined,
    'city': undefined,
    'address': {
      'city': undefined,
      'postNumber': undefined,
      'address': undefined,
    },
    'phoneNumber': undefined,
    'personalNumber': undefined,
    'publicKey': undefined,
    'badges': [],
    'roles': [],
    'societies': [],
    'agreementVersion': '',
  },
  'settings': {
    'appearance': {
      'language': 'swedish',
      'nightMode': false,
      'sort': 'new',
      'display': 'default',
      'newTab': false,
    },
    'notifications': {},
    'privacy': {
      'securityHistory': false,
      'displayName': 'name',
    },
  },
  'profile': {
    'badges': [],
    'colour': undefined,
    'status': undefined,
    'bio': undefined,
    'url': undefined,
  },
  'agora': {
    'followedHypagoras': [],
    'followedUsers': [],
    'starredAgoragrams': [],
    'score': {
      'posts': 0,
      'comments': 0,
      'stars': 0,
      'followers': 0,
    },
  },
  'resetPassword': {},
  'cooldowns': {},
  'connectedApps': {},
  'notifications': [],
  'securityHistory': [],
};

module.exports = {
  name: 'member',
  description: 'Initierar *bli medlem* genom discord',
  aliases: [ 'medlem' ],
  group: 'digitalungdom',
  usage: 'member',
  example: 'member',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( args.length === 0 ) {
      const embed1 = {
        'embed': {
          'title': 'Bli medlem',
          'description': 'Vad kul att du vill bli medlem, du kan även bli medlem hos denna [länk](https://digitalungdom.se/bli-medlem). Följ processen som är förklarad nedan så är du snart medlem.',
          'color': 4086462,
          'fields': [ {
              'name': 'Process',
              'value': 'Du kommer få ett antal frågor av boten (t.ex. vad är din e-post?). Du ska då bara svara boten, utan några kommando. För att avbryta processen tidigit skriv:\n `stop`'
            },
            {
              'name': 'Varför?',
              'value': 'Varför vill man bli medlem? När du är medlem kan du bland annat: lägga upp saker på vårt forum, delta i roliga open-source programmerings projekt och lära känna andra programmerare. Kolla gärna in vår hemsida: [digitalungdom.se](https://digitalungdom.se/bli-medlem)'
            },
            {
              'name': 'GDPR och Stadgar',
              'value': 'Genom att bli medlem accepterar du Digital Ungdoms användarvillkor. Detta innebär:\n att du accepterar att Digital Ungdom lagrar den information du anger på hemsidan. Du godkänner även att du kommer att följa förenings stadgar, som du kan läsa [här](https://digitalungdom.se/stadgar.pdf).'
            },
            {
              'name': 'Hur vi behandlar din data',
              'value': 'Ditt lösenord kommer först och främst hashas, detta innebär att ingen kommer kunna klura ut vad ditt lösenord är. Sedan kommer dina uppgifter lagras i vår databas.'
            }
          ]
        }
      };

      const embedEmail = {
        'embed': {
          'title': 'Vad är din e-post?',
          'description': 'Exempel: `exempel@digitalungdom.se`',
          'color': 4086462
        }
      };

      global.memberProcess[ message.author.id ] = {};
      global.memberProcess[ message.author.id ].case = 'email';
      global.memberProcess[ message.author.id ].user = user;

      message.author.send( embed1 );
      return message.author.send( embedEmail );
    } else if ( args === 'stop' ) {
      delete global.memberProcess[ message.author.id ];

      return message.author.send( 'Processen avbröts.' );
    } else {
      if ( message.channel.type !== 'dm' ) {
        message.delete();
        return message.author.send( 'Du måste skicka dina detaljer här via DM.' );
      }

      let embed;
      let validation;

      switch ( memberProcess[ message.author.id ].case ) {
      case 'email':
        validation = await validate.details.email( args );

        if ( validation.error ) {
          embed = {
            'embed': {
              'title': 'Vad är din e-post?',
              'description': 'Exempel: `exempel@digitalungdom.se`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': validation.error
              } ]
            }
          };

          return message.author.send( embed );
        }

        global.memberProcess[ message.author.id ].user.details.email = args;
        global.memberProcess[ message.author.id ].case = 'username';

        embed = {
          'embed': {
            'title': 'Vilket användarnamn vill du använda?',
            'description': 'Exempel: `exempel`',
            'color': 4086462
          }
        };

        return message.author.send( embed );
      case 'username':
        validation = await validate.details.username( args );

        if ( validation.error ) {
          embed = {
            'embed': {
              'title': 'Vilket användarnamn vill du använda?',
              'description': 'Exempel: `exempel`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': validation.error
              } ]
            }
          };

          return message.author.send( embed );
        }

        global.memberProcess[ message.author.id ].user.details.username = args;
        global.memberProcess[ message.author.id ].case = 'name';

        embed = {
          'embed': {
            'title': 'Vad är ditt fullständiga namn?',
            'description': 'Exempel: `sven svensson`',
            'color': 4086462
          }
        };

        return message.author.send( embed );
      case 'name':
        validation = await validate.details.name( args );

        if ( validation.error ) {
          embed = {
            'embed': {
              'title': 'Vad är ditt fullständiga namn?',
              'description': 'Exempel: `sven svensson`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': validation.error
              } ]
            }
          };

          return message.author.send( embed );
        }
        args = args.toLowerCase().split( ' ' ).filter( n => n ).map( ( s ) => ( [ 'von', 'van', 'de', 'der', 'los', 'ibn', 'd´', 'd\'' ].indexOf( s ) === -1 ) ? s.charAt( 0 ).toUpperCase() + s.substring( 1 ) : s ).join( ' ' );

        global.memberProcess[ message.author.id ].user.details.name = args;
        global.memberProcess[ message.author.id ].case = 'birthdate';

        embed = {
          'embed': {
            'title': 'När är du född?',
            'description': 'Exempel: `2000-04-14`',
            'color': 4086462
          }
        };

        return message.author.send( embed );
      case 'birthdate':
        validation = await validate.details.birthdate( args );

        if ( validation.error ) {
          embed = {
            'embed': {
              'title': 'När är du född?',
              'description': 'Exempel: `2000-04-14`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': validation.error
              } ]
            }
          };

          return message.author.send( embed );
        }

        args = args.split( '-' );
        args = new Date( Date.UTC( args[ 0 ], args[ 1 ] - 1, args[ 2 ] ) );

        global.memberProcess[ message.author.id ].user.details.birthdate = args;
        global.memberProcess[ message.author.id ].case = 'gender';

        embed = {
          'embed': {
            'title': 'Vilket kön har du? (man, kvinna, vill ej uppge eller annat)',
            'description': 'Exempel: `vill ej uppge`',
            'color': 4086462
          }
        };

        return message.author.send( embed );
      case 'gender':
        switch ( args ) {
        case 'man':
          args = 0;
          break;
        case 'kvinna':
          args = 1;
          break;
        case 'vill ej uppge':
          args = 2;
          break;
        case 'annat':
          args = 3;
          break;
        default:
          embed = {
            'embed': {
              'title': 'Vilket kön har du? (man, kvinna, vill ej uppge eller annat)',
              'description': 'Exempel: `vill ej uppge`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': 'Kön måste vara: `man, kvinna, vill ej uppge eller annat`.'
              } ]
            }
          };
          return message.author.send( embed );
        }

        global.memberProcess[ message.author.id ].user.details.gender = args;
        global.memberProcess[ message.author.id ].case = 'password';

        embed = {
          'embed': {
            'title': 'Vad vill du ha för lösenord? Meddelandet tas bort direkt och ditt lösenord hashas samtidigt.',
            'description': 'Exempel: `exempellösenord123`',
            'color': 4086462
          }
        };

        return message.author.send( embed );
      case 'password':
        validation = await validate.details.password( args );

        if ( validation.error ) {
          embed = {
            'embed': {
              'title': 'Vad vill du ha för lösenord? Ditt lösenord hashas direkt.',
              'description': 'Exempel: `exempellösenord123`',
              'color': 16711680,
              'fields': [ {
                'name': 'ERROR',
                'value': validation.error
              } ]
            }
          };

          return message.author.send( embed );
        }

        global.memberProcess[ message.author.id ].user.details.password = args;
        if ( !( await checkDiscordId( message.author.id ) ) ) global.memberProcess[ message.author.id ].user.connectedApps.discord = message.author.id;

        createUser( memberProcess[ message.author.id ].user );
        delete global.memberProcess[ message.author.id ];

        return message.author.send( 'Grattis du är nu medlem! Men innan du kan använda ditt konto måste du verfiera din e-post genom att trycka på den länken vi har skickat till dig.' );
      }
    }
  }
};