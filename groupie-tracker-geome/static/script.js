
        const searchInput = document.getElementById('searchInput');
        const suggestionsContainer = document.getElementById('suggestions');

        searchInput.addEventListener('input', async () => {
            const searchTerm = searchInput.value.trim();
            suggestionsContainer.innerHTML = '';

            if (searchTerm) {
                const response = await fetch(`/autocomplete?search=${searchTerm}`);
                const suggestions = await response.json();

                suggestions.slice(0, 4).forEach(suggestion => {
                    const suggestionItem = document.createElement('div');
                    suggestionItem.textContent = suggestion;
                    suggestionItem.classList.add('suggestion');
                    suggestionsContainer.appendChild(suggestionItem);

                    suggestionItem.addEventListener('click', () => {
                        searchInput.value = suggestion;
                        suggestionsContainer.innerHTML = '';
                    });
                });

                suggestionsContainer.style.display = 'block';
            } else {
                suggestionsContainer.style.display = 'none';
            }
        });

        // Закрывать подсказки при клике вне поля поиска
        document.addEventListener('click', event => {
            if (!event.target.closest('.search')) {
                suggestionsContainer.style.display = 'none';
            }
        });

      document.getElementById("backButton").addEventListener("click", function() {
            // При нажатии кнопки "Вернуться назад" используем JavaScript для возврата на предыдущую страницу
            window.history.back();
        });


// JavaScript для добавления/удаления класса "checked" при изменении состояния чекбоксов
const checkboxes = document.querySelectorAll('input[type="checkbox"]');
checkboxes.forEach(checkbox => {
    checkbox.addEventListener('change', function() {
        const checkboxContainer = this.parentElement;
        if (this.checked) {
            checkboxContainer.classList.add('checked');
        } else {
            checkboxContainer.classList.remove('checked');
        }
    });
});

