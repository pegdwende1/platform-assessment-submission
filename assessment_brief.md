# DevOps Engineer Take-Home Assessment

## Deliverables

Submit the following as part of your assessment.

---

## Platform Proposal

Submit a written proposal in PDF, Word, Google Drive, or Markdown format describing how you would approach the initiative.

Your proposal should explain how you would:

- Establish a predictable release cadence across products.
- Standardize deployment practices to reduce risk where appropriate.
- Enable QA engineers to safely manage routine releases.
- Improve release visibility, rollback capabilities, and observability.
- Build an approach that remains maintainable as the organization grows.

These goals may conflict. Explain how you would prioritize and balance them.

Architecture diagrams are welcome but optional. A concise proposal that demonstrates sound engineering judgment is preferred over a longer document that attempts to address every possible scenario.

---

## Technical Demonstration / POC

Create a simple application and demonstrate a complete CI/CD workflow using the cloud platform and tooling of your choice.

The application itself may be minimal (for example, a static webpage or simple API). The focus of this exercise is the software delivery process rather than the complexity of the application.

Your demonstration should showcase, where practical:

- Source control
- Automated build
- Automated deployment
- Basic validation or testing
- Rollback strategy (or how you would implement one)
- A short README describing your approach and design decisions.

Please include a link to the source repository as part of your submission. Be prepared to discuss your implementation, tradeoffs, and how your approach would scale to a larger engineering organization during the interview.

---

## Your Proposal

### Current State Assessment

Based on the information provided:

- What challenges do you identify?
- What risks concern you most?
- What additional information would you gather before beginning?

### Platform Strategy

Describe your overall vision for the engineering platform.

Discuss:

- The major capabilities your platform would provide.
- The principles that guide your recommendations.
- How you would balance standardization with flexibility across multiple technology stacks.
- Where appropriate, identify technologies you would evaluate and why.

For each major recommendation, discuss:

- Why you selected it.
- Alternatives you considered.
- Assumptions you made.
- Situations where you would recommend a different approach.

### Implementation Roadmap

Assume leadership asks you to deliver this initiative under three different timelines.

Describe how your approach changes if the project must be substantially complete within:

- 3 months
- 6 months
- 12 months

For each timeline, discuss:

- Your highest priorities.
- What work is intentionally deferred.
- What compromises you would and would not accept.

---

## Mid-Project Change: Company Acquisition

Six months into the project, Acme acquires another software company.

The acquisition introduces:

- Approximately 10 additional engineers.
- A separate source control platform.
- Different CI/CD tooling.
- A different cloud provider.
- Different monitoring and observability platforms.
- Several customer-facing applications that must continue operating with minimal disruption.

Leadership has asked you to incorporate the newly acquired products and engineering team into your long-term platform strategy while minimizing disruption to both organizations.

Describe how this changes your original proposal.

Discuss:

- Which recommendations remain unchanged.
- Which priorities shift.
- Whether you would pursue immediate standardization or a phased approach.
- How you would evaluate existing tooling before deciding to replace or retain it.
- How you would integrate two engineering organizations with different technical standards and operational practices.

---

## Building the Platform Team

Although you are the first dedicated DevOps Engineer, several software engineers have expressed interest in infrastructure, automation, observability, and platform engineering.

Describe how you would leverage and mentor these engineers throughout the project while maintaining quality, consistency, and long-term ownership of the platform.

---

## Organizational Adoption

Several senior engineers prefer the existing deployment process because it is familiar and has worked successfully for years.

How would you gain organizational adoption without becoming a bottleneck?

---

## Operational Considerations

Every architectural decision introduces tradeoffs.

Discuss:

- Areas where you would intentionally avoid automation.
- Areas where additional complexity is justified.
- Technical debt you would knowingly accept during the first year.
- How you would determine whether the initiative has been successful.

Consider engineering metrics, operational metrics, business outcomes, and team adoption.
