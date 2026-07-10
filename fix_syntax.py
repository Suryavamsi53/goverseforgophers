import re

with open('ui/templates/pages/practice_questions.html', 'r') as f:
    content = f.read()

replacement = """                    if (labelMatch) {
                        const optLetter = labelMatch[1].replace(')','').replace('.','').toLowerCase();
                        formattedOptions += '<div @click="if (!revealed) { selectedOpt = \\'' + optLetter + '\\'; revealed = true; }" ' +
                            ':class="{' +
                                '\\'border-red-500/50 bg-red-500/10 shadow-[0_0_15px_rgba(239,68,68,0.2)]\\': revealed && selectedOpt === \\'' + optLetter + '\\' && selectedOpt !== correctOpt,' +
                                '\\'border-emerald-500/50 bg-emerald-500/10 shadow-[0_0_15px_rgba(16,185,129,0.2)]\\': revealed && correctOpt === \\'' + optLetter + '\\',' +
                                '\\'border-white/10 bg-white/5 hover:bg-glow/10 hover:border-glow/50 hover:shadow-[0_0_20px_rgba(0,208,132,0.15)]\\': !revealed || (revealed && correctOpt !== \\'' + optLetter + '\\' && selectedOpt !== \\'' + optLetter + '\\')' +
                            '}" ' +
                            'class="flex items-center gap-4 p-4 rounded-xl border transition-all duration-300 cursor-pointer group/opt">' +
                            '<span :class="{' +
                                    '\\'bg-red-500 text-white shadow-[0_0_10px_rgba(239,68,68,0.5)]\\': revealed && selectedOpt === \\'' + optLetter + '\\' && selectedOpt !== correctOpt,' +
                                    '\\'bg-emerald-500 text-white shadow-[0_0_10px_rgba(16,185,129,0.5)]\\': revealed && correctOpt === \\'' + optLetter + '\\',' +
                                    '\\'bg-white/10 text-gray-300 group-hover/opt:bg-glow group-hover/opt:text-black group-hover/opt:shadow-[0_0_10px_rgba(0,208,132,0.5)]\\': !revealed || (revealed && correctOpt !== \\'' + optLetter + '\\' && selectedOpt !== \\'' + optLetter + '\\')' +
                                  '}" ' +
                                  'class="flex items-center justify-center w-8 h-8 rounded-lg font-extrabold text-sm transition-all duration-300 shrink-0">' +
                                optLetter.toUpperCase() +
                            '</span>' +
                            '<span :class="{' +
                                    '\\'text-red-300 font-medium\\': revealed && selectedOpt === \\'' + optLetter + '\\' && selectedOpt !== correctOpt,' +
                                    '\\'text-emerald-300 font-medium\\': revealed && correctOpt === \\'' + optLetter + '\\',' +
                                    '\\'text-gray-300 group-hover/opt:text-white\\': !revealed || (revealed && correctOpt !== \\'' + optLetter + '\\' && selectedOpt !== \\'' + optLetter + '\\')' +
                                  '}" ' +
                                  'class="text-base transition-colors duration-300">' + labelMatch[2].trim() + '</span>' +
                        '</div>';
                    }"""

# Find the start and end of the block to replace
start_marker = "                    if (labelMatch) {"
end_marker = "                    }"

# Find indices
start_idx = content.find(start_marker)
# Find the next closing brace matching the block
end_idx = content.find(end_marker, start_idx) + len(end_marker)

# Replace
new_content = content[:start_idx] + replacement + content[end_idx:]

with open('ui/templates/pages/practice_questions.html', 'w') as f:
    f.write(new_content)

print("Replaced!")
