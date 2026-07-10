(function() {
    let optLetter = 'a';
    let correctOpt = 'b';
    let formattedOptions = '<div @click="if (!revealed) { selectedOpt = \'' + optLetter + '\'; revealed = true; }" ' +
        ':class="{' +
            '\'border-red-500/50 bg-red-500/10 shadow-[0_0_15px_rgba(239,68,68,0.2)]\': revealed && selectedOpt === \'' + optLetter + '\' && selectedOpt !== correctOpt,' +
            '\'border-emerald-500/50 bg-emerald-500/10 shadow-[0_0_15px_rgba(16,185,129,0.2)]\': revealed && correctOpt === \'' + optLetter + '\',' +
            '\'border-white/10 bg-white/5 hover:bg-glow/10 hover:border-glow/50 hover:shadow-[0_0_20px_rgba(0,208,132,0.15)]\': !revealed || (revealed && correctOpt !== \'' + optLetter + '\' && selectedOpt !== \'' + optLetter + '\')' +
        '}" ' +
        'class="flex items-center gap-4 p-4 rounded-xl border transition-all duration-300 cursor-pointer group/opt">' +
        '<span :class="{' +
                '\'bg-red-500 text-white shadow-[0_0_10px_rgba(239,68,68,0.5)]\': revealed && selectedOpt === \'' + optLetter + '\' && selectedOpt !== correctOpt,' +
                '\'bg-emerald-500 text-white shadow-[0_0_10px_rgba(16,185,129,0.5)]\': revealed && correctOpt === \'' + optLetter + '\',' +
                '\'bg-white/10 text-gray-300 group-hover/opt:bg-glow group-hover/opt:text-black group-hover/opt:shadow-[0_0_10px_rgba(0,208,132,0.5)]\': !revealed || (revealed && correctOpt !== \'' + optLetter + '\' && selectedOpt !== \'' + optLetter + '\')' +
              '}" ' +
              'class="flex items-center justify-center w-8 h-8 rounded-lg font-extrabold text-sm transition-all duration-300 shrink-0">' +
            optLetter.toUpperCase() +
        '</span>' +
        '<span :class="{' +
                '\'text-red-300 font-medium\': revealed && selectedOpt === \'' + optLetter + '\' && selectedOpt !== correctOpt,' +
                '\'text-emerald-300 font-medium\': revealed && correctOpt === \'' + optLetter + '\',' +
                '\'text-gray-300 group-hover/opt:text-white\': !revealed || (revealed && correctOpt !== \'' + optLetter + '\' && selectedOpt !== \'' + optLetter + '\')' +
              '}" ' +
              'class="text-base transition-colors duration-300">TEST</span>' +
    '</div>';
})();
