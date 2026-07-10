(function() {
    // mock DOM API to just test for syntax errors
    const optLetter = 'a';
    const correctOpt = 'b';
    let formattedOptions = '<div @click="if (!revealed) { selectedOpt = \'' + optLetter + '\'; revealed = true; }" ' +
        ':class="{' +
            '\'border-red-500/50 bg-red-500/10\': revealed && selectedOpt === \'' + optLetter + '\' && selectedOpt !== correctOpt,' +
            '\'border-emerald-500/50 bg-emerald-500/10\': revealed && correctOpt === \'' + optLetter + '\',' +
            '\'border-white/5 bg-dark-800/40 hover:bg-glow/10 hover:border-glow/50\': !revealed || (revealed && correctOpt !== \'' + optLetter + '\' && selectedOpt !== \'' + optLetter + '\')' +
        '}" ' +
        'class="flex items-center gap-3 p-4 rounded-xl border transition-all cursor-pointer group">' +
        '<span :class="{' +
                '\'bg-red-500/20 text-red-400\': revealed && selectedOpt === \'' + optLetter + '\' && selectedOpt !== correctOpt,' +
                '\'bg-emerald-500/20 text-emerald-400\': revealed && correctOpt === \'' + optLetter + '\',' +
                '\'bg-white/10 text-gray-300 group-hover:bg-glow group-hover:text-white\': !revealed || (revealed && correctOpt !== \'' + optLetter + '\' && selectedOpt !== \'' + optLetter + '\')' +
              '}" ' +
              'class="flex items-center justify-center w-8 h-8 rounded-lg font-bold text-sm transition-colors">' +
            optLetter.toUpperCase() +
        '</span>' +
        '<span class="text-gray-300 text-base group-hover:text-white transition-colors">Test</span>' +
    '</div>';
    console.log(formattedOptions);
})();
