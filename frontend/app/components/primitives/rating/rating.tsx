import { FC, useState } from "react";

type RatingProps = {
  name: string;
  starSize?: number;
};

export const Rating: FC<RatingProps> = ({ name, starSize = 40 }) => {
  const [rating, setRating] = useState(0);

  return (
    <div className="flex flex-row">
      {Array.from({ length: 5 }).map((_, index) => {
        const starValue = index + 1;
        return (
          <label
            key={index}
            htmlFor={`${name}-${starValue}`}
            className="flex items-center"
          >
            <input
              type="radio"
              id={`${name}-${starValue}`}
              name={name}
              value={starValue}
              aria-label={`Rating ${starValue}`}
              onClick={() => setRating(starValue)}
              className="hidden"
            />
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill={starValue <= rating ? "currentColor" : "none"}
              stroke={starValue <= rating ? "currentColor" : "gray"}
              className={`cursor-pointer ${rating >= index + 1 ? "text-yellow-400" : "text-gray-300"}`}
              width={starSize}
              height={starSize}
              onClick={() => setRating(starValue)}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"
              />
            </svg>
          </label>
        );
      })}
      <input type="hidden" name={name} value={rating} />
    </div>
  );
};
