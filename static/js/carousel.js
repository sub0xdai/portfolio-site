class ProjectCarousel {
    constructor(selector) {
        this.container = document.querySelector(selector);
        this.cards = Array.from(this.container.querySelectorAll('.project-card'));
        this.currentIndex = 0;
        this.isTransitioning = false;

        this.initializeControls();
        this.addEventListeners();
        this.enhanceAccessibility();
        this.hideCards();
    }

    initializeControls() {
        this.dotsContainer = document.createElement('div');
        this.dotsContainer.className = 'slider-dots';
        
        this.cards.forEach((_, index) => {
            const dot = document.createElement('div');
            dot.className = 'dot' + (index === 0 ? ' active' : '');
            dot.setAttribute('data-index', index);
            dot.addEventListener('click', () => this.goToSlide(index));
            this.dotsContainer.appendChild(dot);
        });

        this.container.appendChild(this.dotsContainer);
    }

    goToSlide(newIndex) {
        if (this.isTransitioning || newIndex === this.currentIndex) return;

        this.isTransitioning = true;

        this.cards[this.currentIndex].classList.remove('active');
        this.dotsContainer.children[this.currentIndex].classList.remove('active');

        this.cards[newIndex].classList.add('active');
        this.dotsContainer.children[newIndex].classList.add('active');

        this.currentIndex = newIndex;

        setTimeout(() => {
            this.isTransitioning = false;
        }, 300);
    }

    addEventListeners() {
        const nextBtn = document.querySelector('.slider-btn.next');
        const prevBtn = document.querySelector('.slider-btn.prev');

        nextBtn.addEventListener('click', () => this.next());
        prevBtn.addEventListener('click', () => this.previous());

        let touchStartX = 0;
        let touchEndX = 0;

        this.container.addEventListener('touchstart', (e) => {
            touchStartX = e.changedTouches[0].screenX;
        }, { passive: true });

        this.container.addEventListener('touchend', (e) => {
            touchEndX = e.changedTouches[0].screenX;
            this.handleSwipe(touchStartX, touchEndX);
        }, { passive: true });

        document.addEventListener('keydown', (e) => {
            if (e.key === 'ArrowLeft') this.previous();
            if (e.key === 'ArrowRight') this.next();
        });
    }

    handleSwipe(start, end) {
        const threshold = 50;
        const diff = end - start;

        if (Math.abs(diff) > threshold) {
            diff > 0 ? this.previous() : this.next();
        }
    }

    previous() {
        const newIndex = (this.currentIndex - 1 + this.cards.length) % this.cards.length;
        this.goToSlide(newIndex);
    }

    next() {
        const newIndex = (this.currentIndex + 1) % this.cards.length;
        this.goToSlide(newIndex);
    }

    enhanceAccessibility() {
        this.container.setAttribute('role', 'region');
        this.container.setAttribute('aria-label', 'Project Showcase');

        this.cards.forEach((card, index) => {
            card.setAttribute('tabindex', '0');
            card.setAttribute('aria-posinset', index + 1);
            card.setAttribute('aria-setsize', this.cards.length);
        });
    }

    hideCards() {
        this.cards.forEach((card, index) => {
            if (index !== 0) {
                card.classList.add('hidden');
            }
        });
    }
}

class ImprovedProjectCarousel extends ProjectCarousel {
    constructor(selector, autoPlayInterval = 5000) {
        super(selector);
        this.autoPlayInterval = autoPlayInterval;
        this.autoPlayTimer = null;
        this.initAutoPlay();
    }

    initAutoPlay() {
        this.startAutoPlay();
        this.container.addEventListener('mouseenter', () => this.pauseAutoPlay());
        this.container.addEventListener('mouseleave', () => this.startAutoPlay());
    }

    startAutoPlay() {
        this.autoPlayTimer = setInterval(() => {
            this.next();
        }, this.autoPlayInterval);
    }

    pauseAutoPlay() {
        clearInterval(this.autoPlayTimer);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const carousel = new ProjectCarousel('.carousel-container');
});
