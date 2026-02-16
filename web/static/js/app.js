/**
 * LinkBio - Main JavaScript
 * Handles animations, HTMX events, and interactions
 */

// ==========================================================================
// Initialize on DOM Ready
// ==========================================================================
document.addEventListener('DOMContentLoaded', function() {
    initAOS();
    initGSAP();
    initHTMXHandlers();
});

// ==========================================================================
// AOS (Animate on Scroll) Configuration
// ==========================================================================
function initAOS() {
    if (typeof AOS !== 'undefined') {
        AOS.init({
            duration: 600,
            easing: 'ease-out-cubic',
            once: true,
            offset: 50,
            disable: 'mobile' // Disable on mobile for performance
        });
    }
}

// ==========================================================================
// GSAP Animations
// ==========================================================================
function initGSAP() {
    if (typeof gsap === 'undefined') return;

    // Hero section animation
    animateHero();
    
    // Link cards hover effect
    initLinkCardAnimations();
    
    // Button hover effects
    initButtonAnimations();
}

function animateHero() {
    // Navbar fade in
    gsap.from('nav', { 
        opacity: 0, 
        y: -20, 
        duration: 0.8, 
        ease: 'power3.out' 
    });
    
    // Feature cards are animated by AOS (data-aos="fade-up"), no GSAP needed
}

function initLinkCardAnimations() {
    // Profile page link buttons
    const linkButtons = document.querySelectorAll('.link-button');
    if (linkButtons.length > 0) {
        gsap.from('.link-button', {
            opacity: 0,
            y: 30,
            stagger: 0.08,
            duration: 0.6,
            ease: 'power3.out',
            delay: 0.4
        });
    }
}

function initButtonAnimations() {
    // Add hover animations to primary buttons
    document.querySelectorAll('.btn-primary').forEach(function(btn) {
        btn.addEventListener('mouseenter', function() {
            gsap.to(this, { scale: 1.02, duration: 0.2 });
        });
        btn.addEventListener('mouseleave', function() {
            gsap.to(this, { scale: 1, duration: 0.2 });
        });
    });
    
    // Link cards subtle animation
    document.querySelectorAll('.link-card').forEach(function(card) {
        card.addEventListener('mouseenter', function() {
            gsap.to(this, { x: 4, duration: 0.2 });
        });
        card.addEventListener('mouseleave', function() {
            gsap.to(this, { x: 0, duration: 0.2 });
        });
    });
}

// ==========================================================================
// HTMX Event Handlers
// ==========================================================================
function initHTMXHandlers() {
    // Handle error responses
    document.body.addEventListener('htmx:beforeSwap', function(evt) {
        if (evt.detail.xhr.status >= 400) {
            evt.detail.shouldSwap = true;
            const errorHtml = '<div class="error-message rounded-xl px-4 py-3 text-sm animate-slide-in">' 
                + evt.detail.xhr.responseText + '</div>';
            evt.detail.target.innerHTML = errorHtml;
        }
    });
    
    // Animate new elements added via HTMX
    document.body.addEventListener('htmx:afterSwap', function(evt) {
        // Re-init AOS for new elements
        if (typeof AOS !== 'undefined') {
            AOS.refresh();
        }
        
    
    });
    
    // Handle successful form submissions
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.successful) {
            // Could add success toast here
        }
    });
}
