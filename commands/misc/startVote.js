module.exports = {
  name: 'voteStart',
  description: 'Börja en röstning om en fråga',
  aliases: ['br','börjaRösta'],
  group: 'misc',
  usage: 'börjaRösta <id> "<alternativ 1>" "<alternativ 2>"',
  example: '',
  serverOnly: true,
  execute(message, args) {

  	if ( args.length === 0 ) return message.reply( 'du måste ge ett id' );
  	if ( !(args[ 0 ].indexOf('"') === -1) ) return message.reply( 'ditt id kan inte innehålla tecknet "' );
  	if ( args.length === 1) return message.reply( 'du måste minst ha ett alternativ');
    if ( !(args[ 1 ].indexOf('"') === 0)) return message.reply('alternativet måste börja med tecknet "');

    message.delete();
    const id = args[0];
  	const options = args.slice(1).join(' ').split('"').filter((e, i) => i % 2 === 1);
  	voteDic[id] = {options};
    voteDic[id]['users'] = {};
    voteDic[id]['score'] = {};
  	console.log(voteDic);

  },
};