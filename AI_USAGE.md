1. 
I am completing a Junior Software Engineer take-home assignment. This are the requirements:

Build a full-stack calculator application with a React frontend and a backend microservice. The frontend should consume the backend API to perform basic and advanced arithmetic operations. Focus on clean design, maintainable code, and testable architecture.Requirements
Functional
Operations:



Addition, Subtraction, Multiplication, Division

Optional: Exponentiation, Square Root, Percentage



Frontend (React):



Intuitive UI for entering input and displaying results

Input validation and error handling

Responsive design (basic mobile support)



Backend (REST API):



Expose endpoints for calculator operations

Validate input and handle edge cases (division by zero, invalid data)

Return results in JSON format



Non-Functional



Clean, readable, and idiomatic code (frontend and backend)

Unit tests covering key functionality for both layers

Documentation: setup instructions, API usage, and design rationale

Optional: Dockerfile for full-stack deployment



Constraints



Frontend: React (TypeScript preferred)

Backend: Go is perferred



Deliverables



Git repository with frontend and backend code

README with setup instructions, API examples, and design decisions

Unit tests and coverage report

Optional: Dockerfile to run frontend + backend together

My experience is on Python with FastAPI, but I want to build it on Go because that is the company's primary backend language.

I'm relatively new to Go. So first I need you to make me a crash course to learn and understand the essentials of Go to be able to build this app.


2.
I need to define now the architecture we are going to use for the backend, help me with that explaining me the reason of every decision you made. Prioritize idiomatic Go, testability, separation of concerns, and avoiding unnecessary dependencies.

3.
I have built the hole backend infraestracture (without tests yet) I need you to check if there are some important parts I'm missing to have a strong but simple backend for this takeHome task. Keep in mind that I have change some of the decisions, I have made a general function Calculate as an interface that calls the respective function for every operation, this to make simplier the legibility of the code and not to have a hole function with all the calculations and errors managing. Also, I think is simplier if the operands are sent separeted instead of an array, so that's why I make the Request model like that. And I have separate in different files the models and helpers to have the backend more modular and docupled. Before implementing the things you find that are missing, explain me why you think they are missing.

4.
Now that we have it well structured just add the optional operations (Exponentiation, Square Root, Percentage), validate the errors such as fractionary exponent, nevative number in square root, 0^-1. For percentage we are going to define that percentage(a, b) is a% of b. Also make a .md called decisions and write all the important decisions I have make during the implementation like the definition of decisions and other architectural decisions, be concise. Juts write the decision taken: why I took it.
