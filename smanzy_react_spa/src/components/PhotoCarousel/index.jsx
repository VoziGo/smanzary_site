import React, { useState } from 'react';
import styles from './index.module.scss';

// usage example:
// import PhotoCarousel from './PhotoCarousel';
// <PhotoCarousel images={images} />
// const App = () => {
//   const myImages = [
//     'picsum.photos',
//     'picsum.photos',
//     'picsum.photos'
//   ];
//   return <PhotoCarousel images={myImages} />;
// };

const PhotoCarousel = ({ images = [] }) => {
    const [currentIndex, setCurrentIndex] = useState(0);

    const nextSlide = () => {
        if (currentIndex < images.length - 1) {
            setCurrentIndex((prev) => prev + 1);
        }
    };

    const prevSlide = () => {
        if (currentIndex > 0) {
            setCurrentIndex((prev) => prev - 1);
        }
    };

    if (!images.length) return null;

    return (
        <div className={styles.carousel}>
            <div className={styles.carousel__viewport}>
                <div
                    className={styles.carousel__track}
                    style={{ transform: `translateX(-${currentIndex * 100}%)` }}
                >
                    {images.map((url, index) => (
                        <div key={index} className={styles.carousel__slide}>
                            <img src={url} alt={`Slide ${index + 1}`} loading="lazy" />
                        </div>
                    ))}
                </div>
            </div>

            <button
                className={`${styles.carousel__button} ${styles['carousel__button--prev']}`}
                onClick={prevSlide}
                disabled={currentIndex === 0}
                aria-label="Previous slide"
            >
                &#10094;
            </button>

            <button
                className={`${styles.carousel__button} ${styles['carousel__button--next']}`}
                onClick={nextSlide}
                disabled={currentIndex === images.length - 1}
                aria-label="Next slide"
            >
                &#10095;
            </button>
        </div>
    );
};

export default PhotoCarousel;
