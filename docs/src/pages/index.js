import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';

import Heading from '@theme/Heading';
import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title custom__font">
          {siteConfig.tagline}
        </Heading>
        <p className="hero__subtitle">
          Unlock a new level of financial management beyond basic budgeting 
          <br/>
          Get AI-powered insights to questions <i className="custom__font">How much did I spend on dining?</i> or <i className="custom__font">What are my biggest expenses?</i>
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            to="/docs/intro">
            Start Demo
          </Link>
        </div>
      </div>
    </header>
  );
}

function InfoItem({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.benefitSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function HowItWorksItem({Svg, svgDirection, title, description}) {
  return (
    <div className={clsx('row row--align-center margin-top--lg margin-bottom--lg')}>
      <div className={clsx('col col--5 col--offset-1')}>
        {svgDirection == 'left' ? (
          <div className="">
            <Svg className={styles.howItWorksSvg} role="img" />
          </div>
        ) : (
          <div className="">
            <Heading as="h3">{title}</Heading>
            <p>{description}</p>
          </div>
        )}
      </div>
      <div className={clsx('col col--5')}>
        {svgDirection == 'right' ? (
          <div className="">
            <Svg className={styles.howItWorksSvg} role="img" />
          </div>
        ) : (
          <div className="">
            <Heading as="h3">{title}</Heading>
            <p>{description}</p>
          </div>
        )}
      </div>
    </div>
  );
}

const BenefitList = [
  {
    title: 'Easy to Use',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        Docusaurus was designed from the ground up to be easily installed and
        used to get your website up and running quickly.
      </>
    ),
  },
  {
    title: 'Focus on What Matters',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Docusaurus lets you focus on your docs, and we&apos;ll do the chores. Go
        ahead and move your docs into the <code>docs</code> directory.
      </>
    ),
  },
  {
    title: 'Powered by React',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Extend or customize your website layout by reusing React. Docusaurus can
        be extended while reusing the same header and footer.
      </>
    ),
  },
];

const HowItWorksList = [
  {
    title: 'Focus on What Matters',
    svgDirection: 'left',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Docusaurus lets you focus on your docs, and we&apos;ll do the chores. Go
        ahead and move your docs into the <code>docs</code> directory.
      </>
    ),
  },
  {
    title: 'Powered by React',
    svgDirection: 'right',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Extend or customize your website layout by reusing React. Docusaurus can
        be extended while reusing the same header and footer.
      </>
    ),
  },
]

const FAQList = [
  {
    q: 'How can I trust you with my financial data?',
    a: (
      <>
        You OWN your data and security is our top priority. <br /> 
        We are working on adding encryption for keeping it safe and private.
      </>
    )
  },
  {
    q: 'How does the AI-powered querying work?',
    a: (
      <>
        It is powered by state-of-the-art Large Language Model by Meta called <a href="https://www.llama.com">Llama</a>.
      </>
    )
  },
  {
    q: 'Does the app offer notifications or alerts?',
    a: (
      <>
        Not as of today :( But submit a feature request!
      </>
    )
  }
]

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`Welcome`}
      description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <section className={styles.infoItems}>
          <div className="container">
            <Heading as="h1" align="center" className="custom__font">
              Benefits
            </Heading>
            <div className="row">
              {BenefitList.map((props, idx) => (
                <InfoItem key={idx} {...props} />
              ))}
            </div>
          </div>
        </section>
        <section className={styles.infoItems}>
          <div className="container">
            <Heading as="h1" align="center" className="custom__font">
              How It Works?
            </Heading>
            {HowItWorksList.map((props, idx) => (
              <HowItWorksItem key={idx} {...props} />
            ))}
          </div>
        </section>
        <section className={styles.infoItems}>
          <div className="container">
            <Heading as="h1" align="center" className="custom__font">
              FAQ
            </Heading>
            {FAQList.map(({q, a}, idx) => (
              <div key={idx} className={clsx(styles.faqItem, 'padding-top--sm', 'padding-bottom--sm')}>
                — <span style={{'fontSize': '110%', 'fontWeight': 'bold'}}>{q}</span><br/>
                {a}
              </div>
            ))}
            <br />
            <div className={styles.buttons}>
              <Link
                className="button button--primary button--lg"
                to="/docs/intro">
                Start Demo
              </Link>
            </div>
          </div>
        </section>
      </main>
    </Layout>
  );
}
