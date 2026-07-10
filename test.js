(function() {
    const container = document.getElementById('curriculum-container');
    if (!container) return;

    const newContainer = document.createElement('div');
    newContainer.className = "space-y-6 w-full";

    const nodes = Array.from(container.childNodes);
    let currentBlock = null;

    nodes.forEach(node => {
        // Headings
        if (node.tagName && node.tagName.match(/^H[1-6]$/)) {
            if (currentBlock) { newContainer.appendChild(renderBlock(currentBlock)); currentBlock = null; }
            node.className = (node.className || "") + " !mt-12 !mb-6 text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-400 font-extrabold";
            if (node.tagName === 'H1') node.className += " text-4xl border-b border-white/10 pb-4";
            newContainer.appendChild(node.cloneNode(true));
            return;
        }
        
        // HR
        if (node.tagName === 'HR') {
            if (currentBlock) { newContainer.appendChild(renderBlock(currentBlock)); currentBlock = null; }
            const hr = document.createElement('hr');
            hr.className = "border-white/10 my-12";
            newContainer.appendChild(hr);
            return;
        }

        // New Question Block (starts with strong tag like 1.1)
        if (node.tagName === 'P' && node.innerHTML.match(/^<strong>\d+\.\d+<\/strong>/)) {
            if (currentBlock) { newContainer.appendChild(renderBlock(currentBlock)); }
            currentBlock = { type: 'question', nodes: [node.cloneNode(true)] };
            return;
        }

        // Other content
        if (currentBlock) {
            currentBlock.nodes.push(node.cloneNode(true));
        } else {
            // Un-matched content (like intro text)
            if (node.nodeType === 1) { // Element node
                node.className = (node.className || "") + " text-gray-300 text-lg";
            }
            newContainer.appendChild(node.cloneNode(true));
        }
    });

    if (currentBlock) {
        newContainer.appendChild(renderBlock(currentBlock));
    }

    // Replace the original content with our interactive cards
    container.innerHTML = '';
    container.classList.remove('glass-card', 'p-8', 'md:p-12', 'prose'); // Remove default padding and prose
    container.appendChild(newContainer);

    function renderBlock(block) {
        const card = document.createElement('div');
        card.className = "glass-card p-6 md:p-8 relative overflow-hidden transition-all duration-500 hover:shadow-[0_0_40px_rgba(0,208,132,0.1)] hover:border-glow/40 border border-white/10 rounded-2xl group";

        const firstP = block.nodes[0];
        let html = firstP.innerHTML;

        const qContainer = document.createElement('div');
        qContainer.className = "text-gray-100 text-lg md:text-xl font-medium mb-6 leading-relaxed flex flex-col gap-2";

        const parts = html.split(/<br\s*\/?>/i);
        
        let questionText = parts[0];
        let optionsHtml = '';
        let answerText = '';
        let correctOpt = '';

        if (parts.length >= 3 && parts[parts.length-1].includes('Answer:')) {
            // MCQ Format
            questionText = parts[0];
            optionsHtml = parts.slice(1, parts.length - 1).join('<br>');
            answerText = parts[parts.length - 1];
            
            // Extract correct option from answer text, e.g. "Answer: b)"
            const ansMatch = answerText.match(/Answer:\s*([a-f])\)/i);
            if (ansMatch) {
                correctOpt = ansMatch[1].toLowerCase();
            }
            
            // Format options nicely
            const optMatches = optionsHtml.match(/([a-f]\)|[a-f]\.)\s*(.*?)(?=(?:[a-f]\)|[a-f]\.|$))/gi);
            if (optMatches && optMatches.length >= 2) {
                let formattedOptions = '<div class="grid grid-cols-1 md:grid-cols-2 gap-3 mt-4">';
                optMatches.forEach(opt => {
                    const labelMatch = opt.match(/^([a-f]\)|[a-f]\.)\s*(.*)/i);
                    if (labelMatch) {
                        const optLetter = labelMatch[1].replace(')','').replace('.','').toLowerCase();
                        formattedOptions += '<div @click="if (!revealed) { selectedOpt = \'' + optLetter + '\'; revealed = true; }" ' +
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
                                  'class="text-base transition-colors duration-300">' + labelMatch[2].trim() + '</span>' +
                        '</div>';
                    }
                });
                formattedOptions += '</div>';
                optionsHtml = formattedOptions;
            } else {
                optionsHtml = '<div class="glass-panel p-5 rounded-xl font-mono text-sm leading-relaxed mt-4 border border-white/10 text-emerald-300">' + optionsHtml + '</div>';
            }
        } else {
            // Coding Format or No Answer in paragraph
            let answerIndex = parts.findIndex(p => p.includes('Answer:'));
            if (answerIndex !== -1) {
                questionText = parts.slice(0, answerIndex).join('<br>');
                answerText = parts.slice(answerIndex).join('<br>');
            } else {
                questionText = html;
            }
        }

        card.setAttribute('x-data', "{ revealed: false, selectedOpt: null, correctOpt: '" + correctOpt + "' }");

        // Highlight the question number using a pill
        questionText = questionText.replace(/<strong>(\d+\.\d+)<\/strong>/, '<span class="inline-block px-3 py-1 bg-glow/20 text-glow text-sm rounded-full font-bold mr-2 mb-2">$1</span>');

        qContainer.innerHTML = questionText;
        card.appendChild(qContainer);

        if (optionsHtml) {
            const optDiv = document.createElement('div');
            optDiv.innerHTML = optionsHtml;
            card.appendChild(optDiv);
        }

        // Answer Toggle Button
        const btn = document.createElement('button');
        btn.setAttribute('@click', 'revealed = !revealed');
        btn.className = "flex items-center gap-2 px-4 py-2 mt-6 rounded-lg bg-white/5 hover:bg-white/10 text-sm font-semibold text-emerald-400 hover:text-emerald-300 transition-all border border-white/5 hover:border-emerald-500/30 w-fit";
        btn.innerHTML = '<svg x-show="!revealed" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>' +
            '<svg x-show="revealed" style="display:none;" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path></svg>' +
            '<span x-text="revealed ? \'Hide Solution\' : \'View Solution\'"></span>';
        card.appendChild(btn);

        // Answer Wrapper
        const ansWrapper = document.createElement('div');
        ansWrapper.setAttribute('x-show', 'revealed');
        ansWrapper.setAttribute('x-collapse', '');
        ansWrapper.style.display = 'none';
        ansWrapper.className = "mt-4 overflow-hidden";
        
        const ansContent = document.createElement('div');
        ansContent.className = "p-5 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-50 text-base shadow-inner prose prose-invert max-w-none";

        if (answerText) {
            const p = document.createElement('p');
            p.innerHTML = answerText.replace('<strong>Answer:</strong>', '').replace('<strong>Answer:', '<strong>').trim();
            p.className = "font-medium mb-4 text-emerald-300";
            ansContent.appendChild(p);
        }

        for (let i = 1; i < block.nodes.length; i++) {
            const n = block.nodes[i];
            if (n.nodeType === 1) {
                if (n.tagName === 'PRE') {
                    n.className = "bg-dark-900 border border-dark-border rounded-lg p-4 text-sm overflow-x-auto shadow-md";
                }
            }
            ansContent.appendChild(n);
        }

        ansWrapper.appendChild(ansContent);
        card.appendChild(ansWrapper);
        
        return card;
    }
})();
