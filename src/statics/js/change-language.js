const NAMESPACES = [
    'login',
    'forget-password',
    'pre-register',
    'pre-org-register',
    'register',
    'org-register',
    'otp',
    'set-password',
    'profile',
    'confirm',
    'congrats',
];

const RTL_LANGUAGES = ['ar', 'he', 'fa', 'ur'];

async function loadResources(lng) {
    const resources = {};

    for (const ns of NAMESPACES) {
        const res = await fetch(`/statics/locales/${lng}/${ns}.json`).then(r => r.json());
        Object.assign(resources, res);
    }

    return { translation: resources };
}

async function initI18next() {
    const lng = localStorage.getItem('i18nextLng') || 'en';
    const resources = await loadResources(lng);

    i18next.init(
        {
            lng,
            debug: true,
            fallbackLng: 'en',
            resources: { [lng]: resources },
            interpolation: { escapeValue: false },
        },
        () => {
            updateContent();
            updateDirection(lng);
            updateSelectedLanguageUI(lng);
        }
    );
}

function updateContent() {
    //Texts
    document.querySelectorAll('[data-i18n]').forEach(elem => {
        const key = elem.getAttribute('data-i18n');
        elem.textContent = i18next.t(key);
    });

    //Placeholders
    document.querySelectorAll('[data-i18n-placeholder]').forEach(elem => {
        const key = elem.getAttribute('data-i18n-placeholder');
        elem.setAttribute('placeholder', i18next.t(key));
    });
}

function updateDirection(lng) {
    document.body.dir = RTL_LANGUAGES.includes(lng) ? 'rtl' : 'ltr';
}

function updateSelectedLanguageUI(lng) {
    const selectedElements = [
        document.getElementById('desktop-selected'),
        document.getElementById('mobile-selected'),
    ];

    const optionLists = [
        document.querySelectorAll('#desktop-options li'),
        document.querySelectorAll('#mobile-options li'),
    ];

    optionLists.forEach((options, index) => {
        options.forEach(option => {
            option.style.backgroundColor = 'transparent';

            if (option.dataset.value === lng) {
                selectedElements[index].textContent = option.textContent;
                option.style.backgroundColor = '#f9fafb';
            }
        });
    });
}

function changeLanguage(element, lng) {
    const parentElementId = element.parentElement.id;
    localStorage.setItem('i18nextLng', lng);
    updateSelectedLanguageUI(lng);
    initI18next();
    toggleDropdown(parentElementId);
}

function toggleDropdown(id) {
    const options = document.getElementById(id);
    if (options) {
        options.style.display = options.style.display === 'block' ? 'none' : 'block';
    }
}

// Initialize
initI18next();