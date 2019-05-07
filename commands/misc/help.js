/* global include help client */

const helpMainPage = include( 'utils/embeds/helpEmbeds' )[ 'ℹ' ];
const reactions = include( 'utils/embeds/helpEmbeds' ).reactions;

module.exports = {
  name: 'help',
  description: 'Listar alla tillgängliga kommandon.',
  aliases: [ 'commands', 'command', 'hjälp', 'kommando', 'kommandon' ],
  group: 'misc',
  usage: 'help <command>',
  example: 'help play',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( message.channel.type === 'text' && !message.deleted ) message.delete();
    if ( args.length === 0 ) {
      if ( help[ message.author.id ] ) {
        help[ message.author.id ].delete();
        delete global.help[ message.author.id ];
      }

      const helpMessage = await message.author.send( { 'embed': helpMainPage() } );
      global.help[ message.author.id ] = helpMessage;

      for ( const reaction of reactions ) {
        await helpMessage.react( reaction );
      }

      helpMessage.delete( 300000 );
    } else {
      const commandName = args[ 0 ];
      const command = client.commands.get( commandName ) || client.commands.find( cmd => cmd.aliases && cmd.aliases.includes( commandName ) );

      if ( !command ) return message.reply( `det finns inget \`${commandName}\` kommando. Är det fel? Låt oss veta med \`idea <vilket kommando du vill ha>\`` ).then( msg => { msg.delete( 10000 ); } );

      const embed = {
        'title': `**${command.name}**`,
        'description': `*${command.description}*`,
        'color': 4086462,
        'fields': [ {
            'name': 'ANVÄNDNING',
            'value': `>\`${command.usage}\``,
            'inline': true
          },
          {
            'name': 'EXEMPEL',
            'value': `>\`${command.example}\``,
            'inline': true
          },
        ]
      };

      if ( command.aliases.length > 0 ) {
        embed.fields.push( {
          'name': 'ALIAS',
          'value': `\`${command.aliases.join(', ')}\``,
        } );
      }

      return message.reply( { 'embed': embed } );
    }
  },
};